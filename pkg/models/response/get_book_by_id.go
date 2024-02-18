package response

import "pkg/service/pkg/models"

type GetBookById struct {
	Book models.Book `json:"book"`
}
