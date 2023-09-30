package pkgerr_test

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

var _ = Describe("Known", func() {
	It("unwraps into origin error", func() {
		var (
			givenOriginErr = errors.New("origin error")
			givenKnown     = pkgerr.Known{
				Origin: givenOriginErr,
			}
			givenWrappedErr = fmt.Errorf("wrapped: %w", givenKnown)
		)

		Expect(errors.Is(givenWrappedErr, givenOriginErr)).Should(BeTrue())
	})

	It("casts to known error", func() {
		var (
			givenOriginErr = errors.New("origin error")
			givenKnown     = pkgerr.Known{
				Code:   pkgerr.CodeUnauthenticated,
				Simple: givenOriginErr,
			}
			givenWrappedErr = fmt.Errorf("wrapped: %w", givenKnown)
		)

		kerr, ok := pkgerr.AsKnown(givenWrappedErr)
		Expect(ok).Should(BeTrue())
		Expect(kerr).Should(Equal(givenKnown))
	})
})
