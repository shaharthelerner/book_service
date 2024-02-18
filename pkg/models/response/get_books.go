package response

import "pkg/service/pkg/models"

type GetBooks struct {
	Books []models.Book `json:"books"`
}
