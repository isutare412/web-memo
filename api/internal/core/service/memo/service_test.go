package memo

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	Context("validateTags", func() {
		It("passes valid tags", func() {
			var (
				givenTags = []string{"foo", "bar baz"}
			)

			err := validateTags(givenTags)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns error if tag is too long", func() {
			var (
				givenTags = []string{"foo", "bar baz bar baz bar baz bar baz bar baz bar baz"}
			)

			err := validateTags(givenTags)
			Expect(err).Should(HaveOccurred())
		})

		It("returns error if tag is blank", func() {
			var (
				givenTags = []string{"foo", "   \t\n\t  \t"}
			)

			err := validateTags(givenTags)
			Expect(err).Should(HaveOccurred())
		})
	})

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

	DescribeTable("removeReservedTags",
		func(given, expect []string) {
			Expect(removeReservedTags(given)).To(Equal(expect))
		},
		Entry("nil slice", nil, []string{}),
		Entry("empty slice", []string{}, []string{}),
		Entry("contain reserved tags", []string{"foo", tagPublished, "bar"}, []string{"foo", "bar"}),
		Entry("no reserved tags", []string{"foo", "bar"}, []string{"foo", "bar"}),
	)
})
