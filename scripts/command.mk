.PHONY: generate-api
generate-api:
	powershell -ExecutionPolicy Bypass -File .\scripts\generate.ps1

.PHONY: build-mission
build-mission:
	go build -o bin/mission ./cmd/mission

.PHONY: build-drone
build-drone:
	go build -o bin/drone ./cmd/drone

.PHONY: build-telemetry
build-telemetry:
	go build -o bin/telemetry ./cmd/telemetry

.PHONY: build-all
build-all: build-mission build-drone build-telemetry

.PHONY: down
down:
	podman-compose down

.PHONY: cov
cov:
	go test -cover ./... 

.PHONY: mock
mock:
	mockery
