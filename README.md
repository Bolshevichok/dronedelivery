# Drone Delivery

gRPC, Kafka, Redis, PostgreSQL, Docker Compose.

```
docker compose up -d --build
```

Создать миссию:
```
grpcurl -plaintext -d '{\"missions\":[{\"op_id\":1,\"base_id\":1,\"status\":\"created\",\"lat\":55.7558,\"lon\":37.6173,\"alt\":100,\"payload\":1.2}]}' localhost:8080 mission_api.MissionService/UpsertMissions
```

Создать оператора:
```
grpcurl -plaintext -d '{\"operator\":{\"email\":\"operator2@example.com\",\"name\":\"John\\u0020Doe\"}}' localhost:8080 mission_api.MissionService/CreateOperator
```

Создать базу запуска:
```
grpcurl -plaintext -d '{\"launch_base\":{\"name\":\"Base\\u0020Alpha\",\"lat\":55.7558,\"lon\":37.6173,\"alt\":100}}' localhost:8080 mission_api.MissionService/CreateLaunchBase
```

Создать дрона:
```
grpcurl -plaintext -d '{\"drone\":{\"serial\":\"DRONE-002\",\"model\":\"Model\\u0020X\",\"status\":\"available\",\"base_id\":1}}' localhost:8080 mission_api.MissionService/CreateDrone
```


- `missions.created` — миссия создана (Mission Service -> Drone Service)
- `missions.lifecycle` — изменения статуса (Drone Service -> Mission Service)
- `drone.telemetry` — телеметрия (Drone Service -> Telemetry Extractor)
LRANGE telemetry_queue 0 -1