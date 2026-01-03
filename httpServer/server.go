package httpserver

import (
	"fmt"
	"net/http"
	"sync"
)

type message struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

var ChnMapper = sync.Map{}

type ServerConfig struct {
	Host string
	Port string
}

func NewHTTPServer(cfg ServerConfig) *http.Server {

	m := http.NewServeMux()
	registerRoutes(m)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	return &http.Server{
		Handler: m,
		Addr:    addr,
	}

}

func registerRoutes(m *http.ServeMux) {

	// Core Endpoints
	m.Handle("POST /api/send-message", chain(sendMessageHandler(), loggingMiddleware, authMiddleware, validatorMiddleware[sendMessageSchema]()))
	m.Handle("GET /api/events", chain(sseHandler(), loggingMiddleware, authMiddleware))

	// Auth Endpoints
	m.Handle("POST /api/auth/register", chain(registerHandler(), loggingMiddleware, validatorMiddleware[registerUserSchema]()))
	m.Handle("POST /api/auth/login", chain(loginHandler(), loggingMiddleware, validatorMiddleware[loginSchema]()))
	m.Handle("GET /api/auth/logout", chain(logoutHandler(), loggingMiddleware, authMiddleware))
}
