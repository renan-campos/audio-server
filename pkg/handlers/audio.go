package handlers

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/labstack/echo/v4"

	"github.com/renan-campos/audio-server/pkg/storage"
)

const (
	chunkSize  = 1024
	rangeRegex = `bytes=(\d*)-(\d*)`
)

type AudioHandler interface {
	GetAudioFile(c echo.Context) error
	GetAudioHead(c echo.Context) error
}

type audioHandler struct {
	storage    storage.AudioStorageService
	rangeRegex *regexp.Regexp
}

type AudioHandlerSpec struct {
	Storage storage.AudioStorageService
}

func NewAudioHandler(spec AudioHandlerSpec) AudioHandler {
	return &audioHandler{
		storage:    spec.Storage,
		rangeRegex: regexp.MustCompile(rangeRegex),
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

func (h *audioHandler) GetAudioFile(c echo.Context) error {
	const (
		mediaType = "audio/ogg"
	)
	id := c.Param("id")
	// Specify the path to your Ogg sound file
	audioFilePath, err := h.storage.GetAudioFile(id)
	if err != nil {
		return err
	}
	fileSize, err := audioFilePath.FileSize()
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	rangeStart, rangeEnd, err := h.parseRange(
		c.Request().Header.Get("Range"),
		fileSize,
	)

	// Set the appropriate headers for the HTTP response
	header := c.Response().Header()
	header.Set(echo.HeaderContentType, mediaType)
	header.Set(echo.HeaderContentDisposition,
		fmt.Sprintf("attachment; filename=\"%s.ogg\"", id))
	header.Set("Content-Length", fmt.Sprintf("%d", fileSize))
	header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", rangeStart, rangeEnd, fileSize))

	if rangeStart == 0 && rangeEnd == int64(fileSize) {
		return c.File(string(audioFilePath))
	}

	fp, err := os.Open(audioFilePath.FileStr())
	fp.Seek(rangeStart, 0)
	buffer := make([]byte, rangeEnd-rangeStart)
	_, err = fp.Read(buffer)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	c.Response().WriteHeader(http.StatusPartialContent)
	return c.Blob(http.StatusPartialContent, mediaType, buffer)
}

// The Range header abides by the following syntax:
// Range: bytes=0-499
// Range: bytes=-499
// Range: bytes=0-
func (h *audioHandler) parseRange(rangeStr string, fileSize int) (rangeStart int64, rangeEnd int64, err error) {
	fileSize64 := int64(fileSize)
	if rangeStr == "" {
		rangeStart = 0
		rangeEnd = fileSize64
		return
	}
	// FindStringSubmatch returns a slice of the full match and submatches
	matches := h.rangeRegex.FindStringSubmatch(rangeStr)
	if len(matches) != 3 {
		err = fmt.Errorf("invalid Range header")
		return
	}
	if matches[2] == "" {
		rangeEnd = fileSize64
	} else {
		_, err = fmt.Sscanf(matches[2], "%d", &rangeEnd)
		if err != nil {
			return
		}
	}
	if matches[1] == "" {
		rangeStart = fileSize64 - rangeEnd
		rangeEnd = fileSize64
	} else {
		_, err = fmt.Sscanf(matches[1], "%d", &rangeStart)
		if err != nil {
			return
		}
	}
	if rangeStart < 0 ||
		rangeStart > fileSize64 ||
		rangeStart > rangeEnd ||
		rangeEnd > fileSize64 {
		err = fmt.Errorf("invalid Range header")
		return
	}
	return
}
