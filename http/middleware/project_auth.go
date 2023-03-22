package middleware

import (
	"crudly/ctx"
	"crudly/model"
	"crudly/util"
	"crudly/util/result"
	"net/http"

	"github.com/gorilla/mux"
)

type projectAuthInfoGetter interface {
	GetProjectAuthInfo(id model.ProjectId) result.R[model.ProjectAuthInfo]
}

func NewProjectAuth(projectAuthInfoGetter projectAuthInfoGetter) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			projectId := r.Context().Value(ctx.ProjectIdContextKey).(model.ProjectId)

			authInfoResult := projectAuthInfoGetter.GetProjectAuthInfo(projectId)

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

			h.ServeHTTP(w, r)
		})
	}
}
