package memo

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	Context("sortDedupTags", func() {
		It("deduplicates and sort tags", func() {
			var (
				givenTags = []string{
					"guile",
					"ted",
					"Toad",
					"alice",
					"alice",
					"수혁",
					"주은",
					"charlie",
				}
			)

			var (
				wantTags = []string{
					"Toad",
					"alice",
					"charlie",
					"guile",
					"ted",
					"수혁",
					"주은",
				}
			)

			tags := sortDedupTags(givenTags)
			Expect(tags).Should(Equal(wantTags))
		})
	})
})
