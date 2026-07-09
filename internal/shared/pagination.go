package shared

type Pagination struct {
	Size         int   `json:"size"`
	TotalCount   int   `json:"total_count"`
	CurrentPage  int   `json:"current_page"`
	PreviousPage *int  `json:"previous_page,omitempty"`
	NextPage     *int  `json:"next_page,omitempty"`
	TotalPage    int64 `json:"total_page"`
}
