all: bin/audio-server

bin/audio-server: cmd/audio-server/main.go
	go build -o bin/audio-server cmd/audio-server/main.go

clean:
	rm -f bin/audio-server
