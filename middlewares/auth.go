package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/Dryluigi/golang-todo-list/apiHelper/response"
	"github.com/Dryluigi/golang-todo-list/services/auth"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			response.BuildErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		splitted := strings.Split(authHeader, " ")

		if splitted[0] != "Bearer" || splitted[1] == "" {
			response.BuildErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		authData, err := auth.VerifyAndParseAccessToken(splitted[1])

		if err != nil {
			response.BuildErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		newContext := context.WithValue(r.Context(), auth.UserIdContextKey, authData.UserId)

		next.ServeHTTP(w, r.WithContext(newContext))
	})
}
