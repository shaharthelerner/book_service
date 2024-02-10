package response

type GetBooksInventoryResponse struct {
	TotalBooks int `json:"total_books"`
	Authors    int `json:"authors"`
}
