param(
    [Parameter(Position = 0, Mandatory = $True)]
    [ValidateNotNullOrEmpty()]
    [string] $command
)

# Initialization
$ErrorActionPreference = "Stop"

function run() {
    go run main.go
}

function test() {
    go test -v ./...
}

function coverage() {
    New-Item -Force -ItemType directory -Path coverage | Out-Null
    go test -v -covermode=atomic -coverprofile coverage/coverage ./...
    go tool cover -html=coverage/coverage -o coverage/coverage.html
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
    "run",
    "test",
    "coverage",
    "swagger",
    "generate_keypair"
);

$commandToExecute = $commandRegistrations | Where-Object { $_ -eq $command } | Select-Object -First 1
if ($Null -ne $commandToExecute) {
    Invoke-Expression $commandToExecute
} else {
    Write-Host "No profile associated with the phrase: $($command).`n" -ForegroundColor Red
    Write-Output "Available commands:"
    foreach ($entry in $commandRegistrations) { Write-Output "  $($entry.Command)" }
    Write-Output ""
}