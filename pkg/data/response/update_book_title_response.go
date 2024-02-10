package response

type UpdateBookTitleElasticResponse struct {
	Result string                 `json:"result,omitempty"`
	Error  map[string]interface{} `json:"error,omitempty"`
	Status int                    `json:"status,omitempty"`
}
