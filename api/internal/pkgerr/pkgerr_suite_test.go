package pkgerr_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/isutare412/web-memo/api/internal/log"
)

func TestPkgerr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pkgerr Suite")
}

var _ = BeforeSuite(func ()  {
	log.AdaptGinkgo()
})
