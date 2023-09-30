package memo_test

import (
	"testing"

	"github.com/isutare412/web-memo/api/internal/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMemo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Memo Suite")
}

var _ = BeforeSuite(func() {
	log.AdaptGinkgo()
})
