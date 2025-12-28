package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Harichandra-Prasath/Tchat/logging"
	"github.com/go-playground/validator/v10"
)

type middleware func(http.Handler) http.Handler
type ctxKey[T any] struct{}

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
		logging.Logger.Info("", "method", r.Method, "path", r.URL.Path, "time", time.Since(start).Milliseconds())
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func validatorMiddleware[T any]() middleware {

	var key ctxKey[T]

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var body T

			err := json.NewDecoder(r.Body).Decode(&body)
			if err != nil {
				logging.Logger.Info("Error in Unmarshal", "err", err.Error())
				http.Error(w, "Invalid request body", 400)
				return
			}
			if err = validator.New().Struct(body); err != nil {
				logging.Logger.Info("Validation Error", "err", err.Error())
				http.Error(w, err.Error(), 400)
				return
			}
			ctx := context.WithValue(r.Context(), key, body)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
