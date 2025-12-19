# Drone Delivery

gRPC, Kafka, Redis, PostgreSQL, Docker Compose.


```
docker compose up -d --build
```

## CLI (эмуляция UI)


```
go build -o bin/cli ./cmd/cli
```

Создать миссию:

```
./bin/cli create --operator-id 1 --base-id 1 --lat 55.7558 --lon 37.6173 --alt 100 --payload 1.2
```

Смотреть статус и телеметрию стримом:

```
./bin/cli watch <mission-id>
```
## Сборка всех бинарников 

```
go build -o bin/mission ./cmd/mission
go build -o bin/drone ./cmd/drone
go build -o bin/telemetry ./cmd/telemetry
go build -o bin/cli ./cmd/cli
```

## Kafka топики

- `missions.created` — миссия создана (Mission Service -> Drone Service)
- `missions.lifecycle` — изменения статуса (Drone Service -> Mission Service)
- `drone.telemetry` — телеметрия (Drone Service -> Telemetry Extractor)