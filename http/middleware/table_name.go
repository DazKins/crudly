package middleware

import (
	"context"
	"crudly/ctx"
	"crudly/http/dto"
	"net/http"

	"github.com/gorilla/mux"
)

func NewTableName() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)

			tableNameDto := dto.TableNameDto(vars["tableName"])

			tableNameResult := tableNameDto.ToModel()

			if tableNameResult.IsErr() {
				AttachError(w, tableNameResult.UnwrapErr())
				w.WriteHeader(400)
				w.Write([]byte("invalid table name"))
				return
			}

			ctx := context.WithValue(r.Context(), ctx.TableNameContextKey, tableNameResult.Unwrap())

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
