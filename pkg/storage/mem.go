package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

func NewMemAudioStorageService() AudioStorageService {
	return &memAudioStorageService{
		items: map[string]AudioMetadata{},
	}
}

type memAudioStorageService struct {
	items map[string]AudioMetadata
}

func (s *memAudioStorageService) CreateEntry(uuid string) error {
	if _, ok := s.items[uuid]; ok {
		return fmt.Errorf("uuid collision!")
	}
	s.items[uuid] = AudioMetadata{}
	return nil
}

func (s *memAudioStorageService) UpdateMetadata(id string, metadata AudioMetadata) error {
	if _, ok := s.items[id]; !ok {
		return fmt.Errorf("%q not found", id)
	}
	s.items[id] = metadata
	return nil
}

func (s *memAudioStorageService) UploadAudio(id string, file *multipart.FileHeader) error {
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

func (s *memAudioStorageService) ListAudio() (AudioIdList, error) {
	var audioIds []string = []string{}
	for id := range s.items {
		audioIds = append(audioIds, id)
	}
	return AudioIdList{
		Total: len(s.items),
		Items: audioIds,
	}, nil
}

func (s *memAudioStorageService) ListAudioMetadata(id string) (AudioMetadata, error) {
	metadata, ok := s.items[id]
	if !ok {
		return AudioMetadata{}, fmt.Errorf("%q not found", id)
	}
	return metadata, nil
}

func (s *memAudioStorageService) GetAudioFile(id string) (AudioFilePath, error) {
	_, ok := s.items[id]
	if !ok {
		return "", fmt.Errorf("%q not found", id)
	}
	return AudioFilePath(fmt.Sprintf("%s/%s.ogg", uploadDirectory, id)), nil
}
