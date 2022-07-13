package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-profiles-service/internal/users"
	"github.com/rs/zerolog/log"
)

type AuthMiddleware struct {
	UsersClient users.UsersClient
}

func NewAuthMiddleware(usersClient users.UsersClient) AuthMiddleware {
	return AuthMiddleware{
		UsersClient: usersClient,
	}
}

type usernameContextKey int

const UsernameContextKey usernameContextKey = 0

func (h AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const bearerScheme string = "Bearer "

		auth := r.Header.Get("Authorization")
		if len(auth) == 0 {
			auth = r.Header.Get("authorization")
		}

		if !strings.HasPrefix(auth, bearerScheme) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := auth[len(bearerScheme):]

		user, err := h.UsersClient.GetCurrentUser(token)
		if err != nil {
			log.Error().Err(err).Msg("Error getting current user")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		username := user.User.Username
		ctxWithUsername := context.WithValue(r.Context(), UsernameContextKey, username)
		rWithUsername := r.WithContext(ctxWithUsername)
		next.ServeHTTP(w, rWithUsername)
	})
}

func (h AuthMiddleware) AuthenticateOptional(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const bearerScheme string = "Bearer "

		auth := r.Header.Get("Authorization")
		if len(auth) == 0 {
			auth = r.Header.Get("authorization")
		}

		if !strings.HasPrefix(auth, bearerScheme) {
			next.ServeHTTP(w, r)
			return
		}

		token := auth[len(bearerScheme):]

		user, err := h.UsersClient.GetCurrentUser(token)
		if err != nil {
			log.Error().Err(err).Msg("Error getting current user")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		username := user.User.Username
		ctxWithUsername := context.WithValue(r.Context(), UsernameContextKey, username)
		rWithUsername := r.WithContext(ctxWithUsername)
		next.ServeHTTP(w, rWithUsername)
	})
}
