package response

import "pkg/service/pkg/models"

type GetBooksResponse struct {
	Books []models.BookSource `json:"books"` // TODO change to models.Book
}

type GetBooksElasticResponse struct {
	Hits BookHits `json:"hits"`
}

type BookHits struct {
	Hits []BookHit `json:"hits"`
}

type BookHit struct {
	Source models.BookSource `json:"_source"`
}
