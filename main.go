package main

import (
	"github.com/Harichandra-Prasath/Tchat/configs"
	"github.com/Harichandra-Prasath/Tchat/db"
	httpserver "github.com/Harichandra-Prasath/Tchat/httpServer"
	"github.com/Harichandra-Prasath/Tchat/logging"
)

func init() {
	logging.IntialiseLogger()
	configs.InitialiseConfigs()
	err := db.IntialiseDB()
	if err != nil {
		panic(err)
	}
	logging.Logger.Info("DB Initialised")
}

func main() {
	hServer := httpserver.NewHTTPServer(httpserver.ServerConfig{Host: configs.GnCfg.Host, Port: configs.GnCfg.Port})
	logging.Logger.Info("HTTP Server Started", "Host", configs.GnCfg.Host, "Port", configs.GnCfg.Port)
	hServer.ListenAndServe()
}
