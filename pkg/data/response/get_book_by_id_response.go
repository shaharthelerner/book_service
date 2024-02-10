package response

import "pkg/service/pkg/models"

type GetBookByIdResponse struct {
	Book models.Book
}

type GetBookByIdElasticResponse struct {
	Found  bool        `json:"found"`
	Source models.Book `json:"_source"`
}
