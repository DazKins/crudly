package dto

import (
	"crudly/model"
	"crudly/util/optional"
	"crudly/util/result"

	"github.com/google/uuid"
)

type UserIdDto string

func (u UserIdDto) ToModel() result.R[model.UserId] {
	uuid, err := uuid.Parse(string(u))

	if err != nil {
		return result.Errf[model.UserId]("error parsing uuid: %w", err)
	}

	return result.Ok(model.UserId(uuid))
}

func GetUserIdDto(userId model.UserId) UserIdDto {
	return UserIdDto(uuid.UUID(userId).String())
}

type UserDto struct {
	TwitterId *string `json:"twitterId"`
	GoogleId  *string `json:"googleId"`
}

func GetUserDto(user model.User) UserDto {
	return UserDto{
		TwitterId: user.TwitterId.ToPointer(),
		GoogleId:  user.GoogleId.ToPointer(),
	}
}

func (u UserDto) ToModel() result.R[model.User] {
	return result.Ok(model.User{
		TwitterId: optional.FromPointer(u.TwitterId),
		GoogleId:  optional.FromPointer(u.GoogleId),
	})
}
