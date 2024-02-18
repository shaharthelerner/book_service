package response

type GetBooksInventory struct {
	Books   int `json:"books"`
	Authors int `json:"authors"`
}
