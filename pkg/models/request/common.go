package request

type Common struct {
	Username string `json:"username" binding:"required"`
}
