all: bin/audio-server bin/auth-server

bin/audio-server: cmd/audio-server/main.go
	go build -o bin/audio-server cmd/audio-server/main.go

bin/auth-server: cmd/auth-server/main.go
	go build -o bin/auth-server cmd/auth-server/main.go

clean:
	rm -f bin/audio-server bin/auth-server

.PHONY: all bin/audio-server clean
