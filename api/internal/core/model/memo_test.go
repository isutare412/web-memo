package model_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/isutare412/web-memo/api/internal/core/model"
)

var _ = Describe("Memo", func() {
	Context("PaginationParams", func() {
		DescribeTable("GetOrDefault",
			func(params model.PaginationParams, expectPage, expectPageSize int) {
				page, pageSize := params.GetOrDefault()
				Expect(page).To(Equal(expectPage))
				Expect(pageSize).To(Equal(expectPageSize))
			},
			Entry("normal case", model.PaginationParams{
				PageOffset: 123,
				PageSize:   456,
			}, 123, 456),
			Entry("invalid value", model.PaginationParams{
				PageOffset: -1928,
				PageSize:   -2,
			}, 1, 100),
			Entry("zero value", model.PaginationParams{}, 1, 100),
		)
	})
})
