# Drone Delivery

gRPC, Kafka, Redis, PostgreSQL, Docker Compose.

```
docker compose up -d --build
```

Создать миссию:
```
grpcurl -plaintext -d '{\"missions\":[{\"op_id\":1,\"base_id\":1,\"status\":\"created\",\"lat\":55.7558,\"lon\":37.6173,\"alt\":100,\"payload\":1.2}]}' localhost:8080 mission_api.MissionService/UpsertMissions
```

- `missions.created` — миссия создана (Mission Service -> Drone Service)
- `missions.lifecycle` — изменения статуса (Drone Service -> Mission Service)
- `drone.telemetry` — телеметрия (Drone Service -> Telemetry Extractor)
