package enum_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/isutare412/web-memo/api/internal/core/enum"
)

var _ = Describe("Memo", func() {
	Context("SortOrder", func() {
		DescribeTable("GetOrDefault",
			func(given enum.SortOrder, expected enum.SortOrder) {
				Expect(given.GetOrDefault()).To(Equal(expected))
			},
			Entry("SortOrderAsc", enum.SortOrderAsc, enum.SortOrderAsc),
			Entry("SortOrderDesc", enum.SortOrderDesc, enum.SortOrderDesc),
			Entry("zero value", enum.SortOrder(""), enum.SortOrderDesc),
			Entry("invalid value", enum.SortOrder("foo bar"), enum.SortOrderDesc),
		)
	})

	Context("MemoSortKey", func() {
		DescribeTable("GetOrDefault",
			func(given enum.MemoSortKey, expected enum.MemoSortKey) {
				Expect(given.GetOrDefault()).To(Equal(expected))
			},
			Entry("MemoSortKeyCreateTime", enum.MemoSortKeyCreateTime, enum.MemoSortKeyCreateTime),
			Entry("MemoSortKeyUpdateTime", enum.MemoSortKeyUpdateTime, enum.MemoSortKeyUpdateTime),
			Entry("zero value", enum.MemoSortKey(""), enum.MemoSortKeyCreateTime),
			Entry("invalid value", enum.MemoSortKey("foo bar"), enum.MemoSortKeyCreateTime),
		)
	})
})
