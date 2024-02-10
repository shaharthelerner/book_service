package request

import "pkg/service/pkg/models"

type CreateUserActivityRequest struct {
	Username string
	Activity models.UserActivity
}
