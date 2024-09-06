package middlewares

import (
	"context"
	"net/http"

	"git.sr.ht/~rehandaphedar/minv-server/token"
	"github.com/go-chi/render"
)

func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, render.M{
				"error": "Missing token",
			})
			return
		}
		tokenString := cookie.Value

		payload, err := token.VerifyToken(tokenString)
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, render.M{
				"error": "Invalid token",
			})
			return
		}
		ctx := context.WithValue(r.Context(), "authPayload", payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
