package response

type GetBooksResponse struct {
	Hits BookHits `json:"hits"`
}

type BookHits struct {
	Hits []BookHit `json:"hits"`
}

type BookHit struct {
	Source BookSource `json:"_source"`
}

type BookSource struct {
	Title          string  `json:"title"`
	AuthorName     string  `json:"author_name"`
	Price          float64 `json:"price"`
	EbookAvailable bool    `json:"ebook_available"`
	PublishDate    string  `json:"publish_date"`
}
