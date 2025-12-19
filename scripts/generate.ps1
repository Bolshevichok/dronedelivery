# PowerShell script for generating protobuf code on Windows

$projectRoot = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)

$apiDir = Join-Path $projectRoot "api"
$pbDir = Join-Path $projectRoot "internal\pb"
$swaggerDir = Join-Path $pbDir "swagger"

if (!(Test-Path $pbDir)) { New-Item -ItemType Directory -Path $pbDir }
if (!(Test-Path $swaggerDir)) { New-Item -ItemType Directory -Path $swaggerDir }

# Generate gRPC code
protoc -I $apiDir `
  -I (Join-Path $apiDir "google\api") `
  --go_out=$pbDir --go_opt=paths=source_relative `
  --go-grpc_out=$pbDir --go-grpc_opt=paths=source_relative `
  (Join-Path $apiDir "mission\v1\mission.proto")

protoc -I $apiDir `
  -I (Join-Path $apiDir "google\api") `
  --go_out=$pbDir --go_opt=paths=source_relative `
  --go-grpc_out=$pbDir --go-grpc_opt=paths=source_relative `
  (Join-Path $apiDir "models\student_model.proto")

protoc -I $apiDir `
  -I (Join-Path $apiDir "google\api") `
  --go_out=$pbDir --go_opt=paths=source_relative `
  --go-grpc_out=$pbDir --go-grpc_opt=paths=source_relative `
  (Join-Path $apiDir "students_api\students.proto")

# Generate gRPC-Gateway
protoc -I $apiDir `
  -I (Join-Path $apiDir "google\api") `
  --grpc-gateway_out=$pbDir `
  --grpc-gateway_opt=paths=source_relative `
  --grpc-gateway_opt=logtostderr=true `
  (Join-Path $apiDir "mission\v1\mission.proto")

# Generate OpenAPI
protoc -I $apiDir `
  -I (Join-Path $apiDir "google\api") `
  --openapiv2_out=$swaggerDir `
  --openapiv2_opt=logtostderr=true `
  (Join-Path $apiDir "mission\v1\mission.proto")

Write-Host "Protobuf generation completed."