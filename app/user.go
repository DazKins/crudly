package app

import (
	"crudly/model"
	"crudly/util/result"

	"github.com/google/uuid"
)

type userCreator interface {
	CreateUser(id model.UserId, user model.User) error
}

type userFetcher interface {
	FetchUser(id model.UserId) result.R[model.User]
}

type userManager struct {
	userCreator userCreator
	userFetcher userFetcher
}

func NewUserManager(userCreator userCreator, userFetcher userFetcher) userManager {
	return userManager{
		userCreator,
		userFetcher,
	}
}

func (u userManager) CreateUser(user model.User) result.R[model.UserId] {
	id := model.UserId(uuid.New())

	err := u.userCreator.CreateUser(id, user)

	if err != nil {
		return result.Errf[model.UserId]("error creating user: %w", err)
	}

	return result.Ok(id)
}

func (u userManager) GetUser(id model.UserId) result.R[model.User] {
	return u.userFetcher.FetchUser(id)
}
