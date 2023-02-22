package middleware

import (
	"crudly/model"
	"crudly/util"
	"net/http"
)

type projectAuthInfoGetter interface {
	GetProjectAuthInfo(id model.ProjectId) util.Result[model.ProjectAuthInfo]
}

type projectAuth struct {
	projectAuthInfoGetter projectAuthInfoGetter
}

func NewProjectAuth(projectAuthInfoGetter projectAuthInfoGetter) projectAuth {
	return projectAuth{
		projectAuthInfoGetter,
	}
}

func (p projectAuth) Attach(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId := r.Context().Value(util.ProjectIdContextKey).(model.ProjectId)

		authInfoResult := p.projectAuthInfoGetter.GetProjectAuthInfo(projectId)

		if authInfoResult.IsErr() {
			AttachError(w, authInfoResult.UnwrapErr())
			w.WriteHeader(500)
			w.Write([]byte("unexpected error getting project auth details"))
			return
		}

		authInfo := authInfoResult.Unwrap()

		projectKey := r.Header.Get("x-project-key")
		saltedProjectKey := projectKey + authInfo.Salt

		hash := util.StringHash(saltedProjectKey)

		if hash != authInfo.SaltedHash {
			w.WriteHeader(401)
			w.Write([]byte("unauthorized to access project"))
			return
		}

		h(w, r)
	}
}
