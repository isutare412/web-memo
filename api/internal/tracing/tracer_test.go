package tracing_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/isutare412/web-memo/api/internal/tracing"
)

var _ = Describe("Init", func() {
	AfterEach(func() {
		Expect(tracing.Shutdown()).To(Succeed())
	})

	// The service resource is merged with the SDK's default resource, which
	// fails when the two carry different semconv schema URLs.
	It("builds a tracer provider with the SDK default resource", func() {
		err := tracing.Init(tracing.Config{
			Enabled:          true,
			ServiceName:      "webmemo-api-test",
			SamplingRatio:    1,
			OTLPGRPCEndpoint: "127.0.0.1:4317",
		})
		Expect(err).NotTo(HaveOccurred())
	})
})
