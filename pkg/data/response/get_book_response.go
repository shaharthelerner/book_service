package response

type GetBookByIdResponse struct {
	Found  bool       `json:"found"`
	Source BookSource `json:"_source"`
}
