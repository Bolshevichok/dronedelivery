# Drone Delivery

gRPC, Kafka, Redis, PostgreSQL, Docker Compose.

```
docker compose up -d --build
```

Создать миссию:
```
grpcurl -plaintext -d '{\"missions\":[{\"op_id\":1,\"base_id\":1,\"status\":\"created\",\"lat\":55.7558,\"lon\":37.6173,\"alt\":100,\"payload\":1.2}]}' localhost:8080 mission_api.MissionService/UpsertMissions
```
grpcurl -plaintext -d '{"operator":{"email":"operator@example.com","name":"John Doe"}}' localhost:8080 mission_api.MissionService/CreateOperator

grpcurl -plaintext -d '{"launch_base":{"name":"Base Alpha","lat":55.7558,"lon":37.6173,"alt":100}}' localhost:8080 mission_api.MissionService/CreateLaunchBase

grpcurl -plaintext -d '{"drone":{"serial":"DRONE-001","model":"Model X","status":"available","base_id":1}}' localhost:8080 mission_api.MissionService/CreateDrone


- `missions.created` — миссия создана (Mission Service -> Drone Service)
- `missions.lifecycle` — изменения статуса (Drone Service -> Mission Service)
- `drone.telemetry` — телеметрия (Drone Service -> Telemetry Extractor)
LRANGE telemetry_queue 0 -1