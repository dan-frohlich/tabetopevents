package tte

type Paging struct {
	ItemsPerPage       any   `json:"items_per_page"`
	NextPageNumber     int64 `json:"next_page_number"`
	PageNumber         any   `json:"page_number"`
	PreviousPageNumber int64 `json:"previous_page_number"`
	TotalItems         int64 `json:"total_items"`
	TotalPages         int64 `json:"total_pages"`
}
