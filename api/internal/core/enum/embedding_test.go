package enum_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/isutare412/web-memo/api/internal/core/enum"
)

var _ = Describe("Embedding", func() {
	Context("TEIModel", func() {
		DescribeTable("Validate",
			func(given enum.TEIModel, wantErr bool) {
				err := given.Validate()
				if wantErr {
					Expect(err).To(HaveOccurred())
				} else {
					Expect(err).NotTo(HaveOccurred())
				}
			},
			Entry("qwen3 0.6b", enum.TEIModelQwen3Embedding0_6B, false),
			Entry("qwen3 4b", enum.TEIModelQwen3Embedding4B, false),
			Entry("zero value", enum.TEIModel(""), true),
			Entry("invalid value", enum.TEIModel("foo bar"), true),
		)
	})
})
