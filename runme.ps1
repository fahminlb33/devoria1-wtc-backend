function run() {
    go run main.go
}
3
function test() {
    go test
}

function test_cover() {
    go test -race -covermode=atomic -coverprofile=coverage.out
}

function swagger() {
    swag init
}
