package response

import "pkg/service/pkg/models"

type CreateBookResponse struct {
	Book models.Book `json:"book"`
}

type CreateBookElasticResponse struct {
	Result string                 `json:"result,omitempty"`
	Error  map[string]interface{} `json:"error,omitempty"`
	Status int                    `json:"status,omitempty"`
}
