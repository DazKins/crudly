package handler

import (
	"crudly/errs"
	"crudly/http/dto"
	"crudly/http/middleware"
	"crudly/model"
	"crudly/util/result"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type userManager interface {
	CreateUser(user model.User) result.R[model.UserId]
	GetUser(id model.UserId) result.R[model.User]
}

type userHandler struct {
	userManager userManager
}

func NewUserHandler(userManager userManager) userHandler {
	return userHandler{
		userManager,
	}
}

func (u userHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		panic("error reading body")
	}

	var userDto dto.UserDto
	json.Unmarshal(bodyBytes, &userDto)

	userResult := userDto.ToModel()

	if userResult.IsErr() {
		middleware.AttachError(w, userResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid user"))
		return
	}

	user := userResult.Unwrap()

	userIdResult := u.userManager.CreateUser(user)

	if userIdResult.IsErr() {
		middleware.AttachError(w, userIdResult.UnwrapErr())
		w.WriteHeader(500)
		w.Write([]byte("unexpected error creating user"))
		return
	}

	w.Write([]byte(userIdResult.Unwrap().String()))
}

func (u userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userIdDto := dto.EntityIdDto(vars["id"])
	userIdResult := userIdDto.ToModel()

	if userIdResult.IsErr() {
		middleware.AttachError(w, userIdResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid user id"))
		return
	}

	userId := userIdResult.Unwrap()

	userResult := u.userManager.GetUser(model.UserId(userId))

	if userResult.IsErr() {
		err := userResult.UnwrapErr()
		middleware.AttachError(w, err)

		if _, ok := err.(errs.UserNotFoundError); ok {
			w.WriteHeader(404)
			w.Write([]byte("user not found"))
			return
		}

		w.WriteHeader(500)
		w.Write([]byte("unexpected error getting user"))
		return
	}

	user := userResult.Unwrap()

	userDto := dto.GetUserDto(user)

	respBodyBytes, _ := json.Marshal(userDto)

	w.Header().Set("content-type", "application/json")
	w.Write(respBodyBytes)
}
