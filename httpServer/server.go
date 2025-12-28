package httpserver

import (
	"fmt"
	"net/http"
)

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
	m.Handle("POST /api/send-message", chain(sendMessageHandler(), loggingMiddleware, authMiddleware, validatorMiddleware[message]()))
	m.Handle("GET /api/events", chain(sendMessageHandler(), loggingMiddleware, authMiddleware))
}
