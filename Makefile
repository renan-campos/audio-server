all: bin/server

bin/server: cmd/server.go
	go build -o bin/server cmd/server.go

clean:
	rm -f bin/server
