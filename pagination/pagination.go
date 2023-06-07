package pagination

const (
	// DefaultPageSize is the default page size.
	DefaultPageSize = 20
	// MaxPageSize is the max page size.
	MaxPageSize = 99
)

/*
Pagination is a struct for pagination.

NOTICE: Make sure that Pagination.Valid() returns true before querying.
*/
type Pagination struct {
	Page     int64 `form:"page" json:"page" schema:"page" example:"1"`
	PageSize int64 `form:"pageSize" json:"pageSize" schema:"pageSize" example:"20"`
	Count    int64 `form:"count" json:"count" schema:"count" example:"100"`
	MaxPage  int64 `form:"maxPage" json:"maxPage" schema:"maxPage" example:"5"`
}

// Valid determines whether the pagination is valid.
func (pa Pagination) Valid() bool {
	return 1 <= pa.Page && pa.Page <= pa.MaxPage
}

// Offset returns the offset of the query.
//
// NOTICE: Make sure that pa.Valid() returns true.
func (pa Pagination) Offset() int64 {
	return pa.PageSize * (pa.Page - 1)
}

// Limit returns the limit of the query.
//
// NOTICE: Make sure that pa.Valid() returns true.
func (pa Pagination) Limit() int64 {
	remain := pa.Count - pa.Offset()
	if pa.PageSize < remain {
		return pa.PageSize
	}

	return remain
}

// calcMaxPage calculates the max page.
func (pa *Pagination) calcMaxPage() {
	pa.MaxPage = pa.Count / pa.PageSize
	if pa.Count%pa.PageSize != 0 {
		pa.MaxPage++
	}
}

// New creates a new pagination.
func New(page, pageSize, count int64) Pagination {
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	pa := Pagination{Count: count, Page: page, PageSize: pageSize}
	pa.calcMaxPage()

	return pa
}

// Logically logically paginates the slice.
func Logically[T any](s []T, page, pageSize int64) []T {
	start := (page - 1) * pageSize
	if page <= 0 {
		start = 0
	}

	end := start + pageSize
	newSlice := make([]T, 0)

	if page <= 0 {
		return make([]T, 0)
	}
	if start >= int64(len(s)) {
		return make([]T, 0)
	}

	for i := start; i < int64(len(s)) && i < end; i++ {
		newSlice = append(newSlice, s[i])
	}

	return newSlice
}
