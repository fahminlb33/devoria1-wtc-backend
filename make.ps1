param(
    [Parameter(Position = 0, Mandatory = $True)]
    [ValidateNotNullOrEmpty()]
    [string] $command
)

# Initialization
$ErrorActionPreference = "Stop"

function build() {
    go build main.go
}

function run() {
    go run main.go
}

function test() {
    go test -v ./...
}

function coverage() {
    New-Item -Force -ItemType directory -Path coverage | Out-Null
    go test -v -covermode=atomic -coverprofile coverage/coverage.out ./...
    go tool cover -html="coverage/coverage.out" -o coverage/coverage.html
}

function coverage_pretty() {
    New-Item -Force -ItemType directory -Path coverage | Out-Null
    $result = $(gocov test ./... | Out-String)
    $result | gocov-xml > coverage/coverage-gocov.xml
    $result | gocov-html > coverage/coverage-gocov.html
    reportgenerator -reports:coverage/coverage-gocov.xml -targetdir:coverage/reportgen
}

function swagger() {
    swag init
}

function generate_keypair() {
    ssh-keygen -t rsa -b 4096 -m pem -f ./private.pem
    ssh-keygen -f ./private.pem -e -m pem > public.pem
    Remove-Item -Force -Path private.pem.pub | Out-Null
}

# Main Entry Point
$commandRegistrations = @(
    "build",
    "run",
    "test",
    "coverage",
    "coverage_pretty",
    "swagger",
    "generate_keypair"
);

$commandToExecute = $commandRegistrations | Where-Object { $_ -eq $command } | Select-Object -First 1
if ($Null -ne $commandToExecute) {
    Invoke-Expression $commandToExecute
} else {
    Write-Host "No profile associated with the phrase: $($command).`n" -ForegroundColor Red
    Write-Output "Available commands:"
    foreach ($entry in $commandRegistrations) { Write-Output "  $($entry)" }
    Write-Output ""
}