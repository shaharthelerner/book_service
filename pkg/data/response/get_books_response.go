package response

import "pkg/service/pkg/models"

type GetBooksResponse struct {
	Books []models.Book `json:"books"`
}

type GetBooksElasticResponse struct {
	Hits BooksHits `json:"hits"`
}

type BooksHits struct {
	Hits []BookHit `json:"hits"`
}

type BookHit struct {
	Id     string            `json:"_id"`
	Source models.BookSource `json:"_source"`
}
