package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	missionv1 "github.com/Bolshevichok/dronedelivery/internal/pb/mission/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type result struct {
	createLatency time.Duration
	totalLatency  time.Duration
	err           error
}

func main() {
	var (
		target          = flag.String("target", envOr("MISSION_ADDR", "localhost:8080"), "gRPC address of mission-service")
		n               = flag.Int("n", 100, "how many missions to create")
		concurrency     = flag.Int("c", 10, "concurrency (workers)")
		operatorID      = flag.Uint64("operator-id", 1, "operator id")
		baseID          = flag.Uint64("base-id", 1, "launch base id")
		lat             = flag.Float64("lat", 55.7558, "destination latitude base")
		lon             = flag.Float64("lon", 37.6173, "destination longitude base")
		alt             = flag.Float64("alt", 100, "destination altitude")
		payloadKg       = flag.Float64("payload", 1.0, "payload kg base")
		jitterMeters    = flag.Float64("jitter-m", 500, "random jitter around lat/lon in meters")
		waitDelivered   = flag.Bool("wait-delivered", false, "after create, poll until delivered/failed")
		waitTimeout     = flag.Duration("wait-timeout", 45*time.Second, "max wait time per mission when -wait-delivered")
		pollInterval    = flag.Duration("poll-interval", 500*time.Millisecond, "poll interval for GetMission")
		requestTimeout  = flag.Duration("rpc-timeout", 5*time.Second, "timeout per single gRPC call")
		printEvery      = flag.Int("progress", 50, "print progress every N creates")
	)
	flag.Parse()

	if *n <= 0 {
		fmt.Println("n must be > 0")
		os.Exit(2)
	}
	if *concurrency <= 0 {
		fmt.Println("c must be > 0")
		os.Exit(2)
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	conn, err := grpc.Dial(*target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("failed to dial %s: %v\n", *target, err)
		os.Exit(1)
	}
	defer conn.Close()

	client := missionv1.NewMissionServiceClient(conn)

	jobs := make(chan int)
	results := make(chan result, *n)

	var created uint64
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(*concurrency)
	for w := 0; w < *concurrency; w++ {
		go func() {
			defer wg.Done()
			for range jobs {
				res := runOne(rnd, client, *operatorID, *baseID, *lat, *lon, *alt, *payloadKg, *jitterMeters, *waitDelivered, *waitTimeout, *pollInterval, *requestTimeout)
				results <- res
				cur := atomic.AddUint64(&created, 1)
				if *printEvery > 0 && cur%uint64(*printEvery) == 0 {
					elapsed := time.Since(start).Truncate(time.Millisecond)
					fmt.Printf("progress: %d/%d in %s\n", cur, *n, elapsed)
				}
			}
		}()
	}

	go func() {
		for i := 0; i < *n; i++ {
			jobs <- i
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	var (
		success        int
		failed         int
		createLatencies []time.Duration
		totalLatencies  []time.Duration
	)

	for r := range results {
		if r.err != nil {
			failed++
			continue
		}
		success++
		createLatencies = append(createLatencies, r.createLatency)
		totalLatencies = append(totalLatencies, r.totalLatency)
	}

	elapsed := time.Since(start)

	fmt.Println("--- loadtest summary ---")
	fmt.Printf("target: %s\n", *target)
	fmt.Printf("missions: %d, concurrency: %d\n", *n, *concurrency)
	fmt.Printf("success: %d, failed: %d\n", success, failed)
	fmt.Printf("elapsed: %s, rate: %.2f req/s\n", elapsed.Truncate(time.Millisecond), float64(success+failed)/elapsed.Seconds())

	printLat("create", createLatencies)
	if *waitDelivered {
		printLat("total", totalLatencies)
	}
}

func runOne(rnd *rand.Rand, client missionv1.MissionServiceClient, operatorID, baseID uint64, latBase, lonBase, altBase, payloadBase, jitterM float64, waitDelivered bool, waitTimeout, pollInterval, rpcTimeout time.Duration) result {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
	defer cancel()

	lat, lon := jitterLatLon(rnd, latBase, lonBase, jitterM)
	payload := payloadBase * (0.7 + rnd.Float64()*0.6) // +-30%

	createStart := time.Now()
	createResp, err := client.CreateMission(ctx, &missionv1.CreateMissionRequest{
		OperatorId:     operatorID,
		LaunchBaseId:   baseID,
		DestinationLat: lat,
		DestinationLon: lon,
		DestinationAlt: altBase,
		PayloadKg:      payload,
	})
	createLatency := time.Since(createStart)
	if err != nil {
		return result{err: err}
	}

	if !waitDelivered {
		return result{createLatency: createLatency, totalLatency: time.Since(start)}
	}

	deadline := time.Now().Add(waitTimeout)
	for time.Now().Before(deadline) {
		callCtx, callCancel := context.WithTimeout(context.Background(), rpcTimeout)
		resp, err := client.GetMission(callCtx, &missionv1.GetMissionRequest{MissionId: createResp.MissionId})
		callCancel()
		if err == nil && resp.Mission != nil {
			s := resp.Mission.Status
			if s == "delivered" || s == "failed" {
				return result{createLatency: createLatency, totalLatency: time.Since(start)}
			}
		}
		time.Sleep(pollInterval)
	}

	return result{err: fmt.Errorf("timeout waiting mission=%d", createResp.MissionId)}
}

func jitterLatLon(rnd *rand.Rand, lat, lon, jitterM float64) (float64, float64) {
	// Очень грубо: 1 deg lat ~ 111_320 m
	const metersPerDegLat = 111_320.0
	metersPerDegLon := metersPerDegLat * cosDeg(lat)

	dLat := (rnd.Float64()*2 - 1) * jitterM / metersPerDegLat
	dLon := (rnd.Float64()*2 - 1) * jitterM / metersPerDegLon

	return lat + dLat, lon + dLon
}

func cosDeg(deg float64) float64 {
	return math.Cos(deg * (math.Pi / 180.0))
}

func printLat(name string, vals []time.Duration) {
	if len(vals) == 0 {
		fmt.Printf("%s latency: n/a\n", name)
		return
	}
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	p50 := percentile(vals, 50)
	p95 := percentile(vals, 95)
	p99 := percentile(vals, 99)
	min := vals[0]
	max := vals[len(vals)-1]
	fmt.Printf("%s latency: min=%s p50=%s p95=%s p99=%s max=%s\n",
		name,
		min.Truncate(time.Millisecond),
		p50.Truncate(time.Millisecond),
		p95.Truncate(time.Millisecond),
		p99.Truncate(time.Millisecond),
		max.Truncate(time.Millisecond),
	)
}

func percentile(vals []time.Duration, p int) time.Duration {
	if len(vals) == 0 {
		return 0
	}
	if p <= 0 {
		return vals[0]
	}
	if p >= 100 {
		return vals[len(vals)-1]
	}
	idx := int(float64(len(vals)-1) * (float64(p) / 100.0))
	return vals[idx]
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
