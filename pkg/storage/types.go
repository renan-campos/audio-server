package storage

import (
	"mime/multipart"
	"os"
)

type AudioStorageService interface {
	CreateEntry(uuid string) error
	UpdateMetadata(id string, metadata AudioMetadata) error
	UploadAudio(id string, file *multipart.FileHeader) error
	ListAudio() (AudioIdList, error)
	ListAudioMetadata(id string) (AudioMetadata, error)
	GetAudioFile(id string) (AudioFilePath, error)
}

type AudioMetadata struct {
	Name string
}

type AudioIdList struct {
	Total int
	Items []string
}

type AudioFilePath string

func (a AudioFilePath) FileSize() (int, error) {
	fileInfo, err := os.Stat(string(a))
	if err != nil {
		return 0, err
	}
	return int(fileInfo.Size()), nil
}

func (a AudioFilePath) FileStr() string {
	return string(a)
}
