package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"go.etcd.io/etcd/client/v3"
)

func NewEtcdAudioStorageService() AudioStorageService {
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

func (s *etcdStorageService) CreateEntry(uuid string) error {
	var metadata AudioMetadata
	_, err := s.client.KV.Put(s.newContext(), "audio."+uuid, s.marshalData(metadata))
	return err
}

func (s *etcdStorageService) UpdateMetadata(id string, metadata AudioMetadata) error {
	_, err := s.client.KV.Put(s.newContext(), "audio."+id, s.marshalData(metadata))
	return err
}

func (s *etcdStorageService) UploadAudio(id string, file *multipart.FileHeader) error {
	// Open the file for writing
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Specify the path where you want to save the uploaded file
	dstPath := fmt.Sprintf("%s/%s.ogg", uploadDirectory, id) // Change the path as needed

	// Create or open the destination file
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the uploaded file data to the destination file
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}

func (s *etcdStorageService) ListAudio() (AudioIdList, error) {
	audioIdList := AudioIdList{
		Items: []string{},
		Total: 0,
	}
	resp, err := s.client.KV.Get(s.newContext(), "audio.", clientv3.WithPrefix())
	if err != nil {
		return audioIdList, err
	}

	var audioIds []string = []string{}
	for _, audioElement := range resp.Kvs {
		audioIds = append(audioIds, s.pruneKey(string(audioElement.Key)))
	}
	return AudioIdList{
		Total: len(audioIds),
		Items: audioIds,
	}, nil
}

func (s *etcdStorageService) ListAudioMetadata(id string) (AudioMetadata, error) {
	resp, err := s.client.KV.Get(s.newContext(), "audio."+id)
	if err != nil {
		return AudioMetadata{}, err
	}
	if len(resp.Kvs) < 1 {
		return AudioMetadata{}, fmt.Errorf("Expected at least one element")
	}
	return s.unmarshalData(resp.Kvs[0].Value)
}

func (s *etcdStorageService) GetAudioFile(id string) (AudioFilePath, error) {
	_, err := s.ListAudioMetadata(id)
	if err != nil {
		return "", fmt.Errorf("%q not found", id)
	}
	return AudioFilePath(fmt.Sprintf("%s/%s.ogg", uploadDirectory, id)), nil
}

func (s *etcdStorageService) newContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}

func (s *etcdStorageService) unmarshalData(data []byte) (AudioMetadata, error) {
	var audioMetadata AudioMetadata
	err := json.Unmarshal(data, &audioMetadata)
	return audioMetadata, err
}

func (s *etcdStorageService) marshalData(metadata AudioMetadata) string {
	data, err := json.Marshal(metadata)
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
