package model

type SortDirection string

const (
	SortDirectionUnspecified SortDirection = ""
	SortDirectionAsc         SortDirection = "asc"
	SortDirectionDesc        SortDirection = "desc"
)

type QueryOption struct {
	PageOffset int
	PageSize   int
	Direction  SortDirection
}
