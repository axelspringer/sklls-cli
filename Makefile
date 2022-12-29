build:
	env GOOS=darwin go build -ldflags "-s -w" -o ./bin/osx-intel/sklls ./cmd/sklls/*.go
	env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/linux/sklls ./cmd/sklls/*.go
start:
	go run ./cmd/sklls/sklls.go
install: build install-locally
install-locally:
	ln -s "$(PWD)/bin/osx-intel/sklls" /usr/local/bin/sklls
test:
	go test `go list ./... | grep -v cmd/playground`
	
benchmark:
	go test -bench . `go list ./... | grep -v cmd/playground`
