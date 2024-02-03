package enum

type SortOrder string

const (
	SortOrderUnspecified SortOrder = ""
	SortOrderAsc         SortOrder = "asc"
	SortOrderDesc        SortOrder = "desc"
)

func (o SortOrder) GetOrDefault() SortOrder {
	switch o {
	case SortOrderAsc:
		fallthrough
	case SortOrderDesc:
		return o
	default:
		return SortOrderDesc
	}
}

type MemoSortKey string

const (
	MemoSortKeyUnspecified MemoSortKey = ""
	MemoSortKeyCreateTime  MemoSortKey = "createTime"
	MemoSortKeyUpdateTime  MemoSortKey = "updateTime"
)

func (k MemoSortKey) GetOrDefault() MemoSortKey {
	switch k {
	case MemoSortKeyCreateTime:
		fallthrough
	case MemoSortKeyUpdateTime:
		return k
	default:
		return MemoSortKeyCreateTime
	}
}
