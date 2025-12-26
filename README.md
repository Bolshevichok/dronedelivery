# Drone Delivery

gRPC, Kafka, Redis, PostgreSQL, Docker Compose.

```
docker compose up -d --build
```

Создать миссию:
```
& "grpcurl.exe" -plaintext -d '{"missions":[{"operator_id":1,"launch_base_id":1,"destination_lat":55.7558,"destination_lon":37.6173,"destination_alt":100,"payload_kg":1.2}]}' localhost:8080 mission_api.MissionService/UpsertMissions
```


- `missions.created` — миссия создана (Mission Service -> Drone Service)
- `missions.lifecycle` — изменения статуса (Drone Service -> Mission Service)
- `drone.telemetry` — телеметрия (Drone Service -> Telemetry Extractor)
