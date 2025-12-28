package main

import (
	"os"

	httpserver "github.com/Harichandra-Prasath/Tchat/httpServer"
	"github.com/Harichandra-Prasath/Tchat/logging"
)

func init() {
	logging.IntialiseLogger()
}

func main() {
	hServer := httpserver.NewHTTPServer(httpserver.ServerConfig{Host: os.Getenv("HOST"), Port: os.Getenv("PORT")})
	logging.Logger.Info("HTTP Server Started", "Host", os.Getenv("HOST"), "Port", os.Getenv("PORT"))
	hServer.ListenAndServe()
}
