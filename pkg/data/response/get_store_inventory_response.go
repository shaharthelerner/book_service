package response

type GetBooksInventoryResponse struct {
	Books   int `json:"books"`
	Authors int `json:"authors"`
}

type GetInventoryElasticResponse struct {
	BookHits     Hits              `json:"hits"`
	Aggregations StoreAggregations `json:"aggregations"`
}

type Hits struct {
	Total InventoryResultSummary `json:"total"`
}

type InventoryResultSummary struct {
	Value int `json:"value"`
}

type StoreAggregations struct {
	UniqueAuthors UniqueAuthorsAggregation `json:"unique_authors"`
}

type UniqueAuthorsAggregation struct {
	Value int `json:"value"`
}
