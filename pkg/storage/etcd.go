package storage

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewEtcdMovieStorageService() MovieStorageService {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	return &etcdStorageService{
		client: etcdClient,
	}
}

type etcdStorageService struct {
	client *clientv3.Client
}

func (s *etcdStorageService) CreateEntry(movieName string) error {
	metadata := MovieData{
		Name: movieName,
	}
	_, err := s.client.KV.Put(s.newContext(), "movie."+movieName, s.marshalData(metadata))
	return err
}

func (s *etcdStorageService) ListMovies() (MovieList, error) {
	resp, err := s.client.KV.Get(s.newContext(), "movie.", clientv3.WithPrefix())
	if err != nil {
		return MovieList{}, err
	}

	var movies []string = []string{}
	for _, movieElement := range resp.Kvs {
		movies = append(movies, s.pruneKey(string(movieElement.Key)))
	}
	return MovieList{
		Total: len(movies),
		Items: movies,
	}, nil
}

func (s *etcdStorageService) newContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}

func (s *etcdStorageService) unmarshalData(data []byte) (MovieData, error) {
	var movieData MovieData
	err := json.Unmarshal(data, &movieData)
	return movieData, err
}

func (s *etcdStorageService) marshalData(movieData MovieData) string {
	data, err := json.Marshal(movieData)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (s *etcdStorageService) pruneKey(input string) string {
	// Find the position of the period
	periodIndex := strings.Index(input, ".")

	// Check if a period is found
	if periodIndex == -1 {
		return input
	}

	// Extract everything after the period
	return input[periodIndex+1:]
}
