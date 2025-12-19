# Drone Delivery System
``` 
docker compose up -d --build
```
to test:
``` 
./bin/cli create --operator-id 1 --base-id 1 --lat 55.7558 --lon 37.6173 --alt 100 --payload 5
```

- Kafka topics: missions.created, missions.lifecycle, drone.telemetry