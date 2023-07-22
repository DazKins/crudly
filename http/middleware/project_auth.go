package middleware

import (
	"crudly/ctx"
	"crudly/errs"
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
			if r.Context().Value(AdminContextKey) != nil {
				h.ServeHTTP(w, r)
				return
			}

			projectId, ok := r.Context().Value(ctx.ProjectIdContextKey).(model.ProjectId)

			if !ok {
				w.WriteHeader(400)
				w.Write([]byte("project id header is not present"))
				return
			}

			authInfoResult := projectAuthInfoGetter.GetProjectAuthInfo(projectId)

			if authInfoResult.IsErr() {
				err := authInfoResult.UnwrapErr()

				if _, ok := err.(errs.ProjectNotFoundError); ok {
					w.WriteHeader(404)
					w.Write([]byte("project not found"))
					return
				}

				AttachError(w, err)
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
