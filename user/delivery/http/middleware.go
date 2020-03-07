package http

import (
	"context"
	"dbProject/user"
	"github.com/dimfeld/httptreemux"
	"net/http"
)

func AuthMiddleware(f httptreemux.HandlerFunc, uc user.UseCase) httptreemux.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		cookie, err := r.Cookie("Auth")
		ctx := r.Context()
		if err != nil {
			r = r.WithContext(context.WithValue(ctx, "user", nil))
			f(w, r, ps)
			return
		}

		u, err := uc.CheckAuth(r.Context(), cookie.Value)
		if err != nil {
			r = r.WithContext(context.WithValue(ctx, "user", nil))
		} else {
			r = r.WithContext(context.WithValue(ctx, "user", u))
		}

		f(w, r, ps)
	}
}