package lifecycle_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("v0 Lifecycle Operations", func() {
	Context("adding a new audio entry", func() {
		It("generates a unique id", func() {
		})
	})
	Context("renaming the audio entry", func() {
		It("updates the metadata for the audio entry", func() {
		})
	})
	Context("upload file to the audio entry", func() {
		It("successfully uploads the file", func() {
		})
	})
})

type Client interface {
	ListAudio()
}

func NewClient(backendAddr string) Client {
	return &clientImpl{
		backendAddr: backendAddr,
	}
}

type clientImpl struct {
	backendAddr string
}

func (c *clientImpl) ListAudio() {
	resp, err := http.Get(
		fmt.Sprintf("%s/v0/audio", c.backendAddr),
	)
	Expect(err).ToNot(HaveOccurred())
}
