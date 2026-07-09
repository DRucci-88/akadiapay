package shared

import "math"

type Page[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

func NewPageEmpty[T any](
	pageable *Pageable,
) *Page[T] {
	pageable.Normalize()
	return &Page[T]{
		Data: make([]T, 0),
		Pagination: Pagination{
			Size:         pageable.Size,
			TotalCount:   0,
			CurrentPage:  pageable.Page,
			PreviousPage: nil,
			NextPage:     nil,
			TotalPage:    0,
		},
	}
}

func NewPage[T any](
	pageable *Pageable,
	totalCount int64,
	data []T,
) *Page[T] {

	pageable.Normalize()

	totalPage := int(math.Ceil(float64(totalCount) / float64(pageable.Size)))

	var previousPage *int
	if pageable.Page > 1 {
		p := pageable.Page - 1
		previousPage = &p
	}

	var nextPage *int
	if pageable.Page < totalPage {
		n := pageable.Page + 1
		nextPage = &n
	}

	return &Page[T]{
		Data: data,
		Pagination: Pagination{
			Size:         pageable.Size,
			TotalCount:   int(totalCount),
			CurrentPage:  pageable.Page,
			PreviousPage: previousPage,
			NextPage:     nextPage,
			TotalPage:    int64(totalPage),
		},
	}
}
