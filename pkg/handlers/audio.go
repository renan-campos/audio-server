package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/renan-campos/audio-server/pkg/storage"
)

const (
	chunkSize = 1024
)

type AudioHandler interface {
	GetAudioHead(c echo.Context) error
}

type audioHandler struct {
	storage storage.AudioStorageService
}

type AudioHandlerSpec struct {
	Storage storage.AudioStorageService
}

func NewAudioHandler(spec AudioHandlerSpec) AudioHandler {
	return &audioHandler{
		storage: spec.Storage,
	}
}

func (h *audioHandler) GetAudioHead(c echo.Context) error {
	id := c.Param("id")
	audioFilePath, err := h.storage.GetAudioFile(id)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	fileSize, err := audioFilePath.FileSize()
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	response := c.Response()
	header := response.Header()
	header.Set("X-Chunk-Size", fmt.Sprintf("%d", chunkSize))
	header.Set("Content-Length", fmt.Sprintf("%d", fileSize))
	c.Response().WriteHeader(http.StatusOK)
	return nil
}
