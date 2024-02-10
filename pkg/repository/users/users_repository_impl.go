package users_repository

type UsersRepositoryImpl struct {
}

func NewUsersRepositoryImpl() UsersRepository {
	return &UsersRepositoryImpl{}
}

func (u UsersRepositoryImpl) SaveActivity() error {
	//TODO implement me
	panic("implement me")
}

func (u UsersRepositoryImpl) GetActivities() error {
	//TODO implement me
	panic("implement me")
}
