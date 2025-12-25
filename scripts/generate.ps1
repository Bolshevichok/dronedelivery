# PowerShell script for generating protobuf code on Windows

$projectRoot = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)

$apiDir = Join-Path $projectRoot "api"
$pbDir = Join-Path $projectRoot "internal\pb"

if (!(Test-Path $pbDir)) { New-Item -ItemType Directory -Path $pbDir }

# Add Go bin to PATH
$goPath = & go env GOPATH
$goBin = Join-Path $goPath "bin"
$env:PATH = "$goBin;$env:PATH"

# Generate gRPC code for models
protoc -I $apiDir `
  --go_out=$pbDir --go_opt=paths=source_relative `
  (Join-Path $apiDir "models\models.proto")

# Generate gRPC code for mission API
protoc -I $apiDir `
  --go_out=$pbDir --go_opt=paths=source_relative `
  --go-grpc_out=$pbDir --go-grpc_opt=paths=source_relative `
  (Join-Path $apiDir "mission_api\mission.proto")

Write-Host "Protobuf generation completed."