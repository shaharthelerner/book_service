package controller

import "pkg/service/pkg/data/request"

func (lc *LibraryController) saveUserAction(username string, method string, route string) error {
	ua := request.CreateUserActivityRequest{
		Username: username,
		Method:   method,
		Route:    route,
	}
	return lc.usersService.SaveUserAction(ua)
}
