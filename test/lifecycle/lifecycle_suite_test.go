package lifecycle_test

import (
	"flag"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var backendAddr *string

func TestTest(t *testing.T) {
	backendAddr = flag.String("backend-addr", "127.0.0.1:1323", "The address the API calls will be made to")
	flag.Parse()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}
