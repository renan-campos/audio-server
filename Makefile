all: bin/movie-server

bin/movie-server: cmd/movie-server/main.go
	go build -o bin/movie-server cmd/movie-server/main.go

clean:
	rm -f bin/movie-server bin/auth-server

.PHONY: all bin/movie-server clean
