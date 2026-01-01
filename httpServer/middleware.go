package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Harichandra-Prasath/Tchat/db"
	"github.com/Harichandra-Prasath/Tchat/logging"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type middleware func(http.Handler) http.Handler
type ctxKey[T any] struct{}

type userIdKey uuid.UUID

func chain(h http.Handler, m ...middleware) http.Handler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logging.Logger.Info("Request Info", "method", r.Method, "path", r.URL.Path, "time", time.Since(start).Milliseconds())
	})
}

func authMiddleware(next http.Handler) http.Handler {

	var key userIdKey

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			logging.Logger.Info("Authorization header missing")
			http.Error(w, "Auth header missing", 401)
			return
		}

		if !strings.HasPrefix(tokenHeader, "Bearer") {
			logging.Logger.Info("Invalid Auth Header. Bearer Missing")
			http.Error(w, "Invalid Auth Header. Bearer Missing", 401)
			return
		}

		token := strings.Split(tokenHeader, " ")[1]

		userId, err := db.GetSession(token)
		if err != nil {
			if errors.Is(err, db.SessionDoesNotExistsError) {
				http.Error(w, "Session doesn't exist. Please login again", 401)
				return
			}

			logging.Logger.Error("Auth Failed", "err", err.Error())
			http.Error(w, "Auth Failed", 500)
			return
		}

		ctx := context.WithValue(r.Context(), key, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validatorMiddleware[T any]() middleware {

	var key ctxKey[T]

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var body T

			err := json.NewDecoder(r.Body).Decode(&body)
			if err != nil {
				logging.Logger.Error("Error in Unmarshal", "err", err.Error())
				http.Error(w, "Invalid request body", 400)
				return
			}

			err = validator.New().Struct(body)
			if err != nil {
				logging.Logger.Error("Validation Error", "err", err.Error())
				http.Error(w, err.Error(), 400)
				return
			}
			ctx := context.WithValue(r.Context(), key, body)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
