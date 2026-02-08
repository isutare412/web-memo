package embedding

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/isutare412/web-memo/api/internal/log"
)

func TestEmbedding(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Embedding Suite")
}

var _ = BeforeSuite(func() {
	log.AdaptGinkgo()
})
