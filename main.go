package main

import (
	"github.com/Harichandra-Prasath/Tchat/broker"
	"github.com/Harichandra-Prasath/Tchat/configs"
	"github.com/Harichandra-Prasath/Tchat/db"
	httpserver "github.com/Harichandra-Prasath/Tchat/httpServer"
	"github.com/Harichandra-Prasath/Tchat/logging"
)

func init() {
	logging.IntialiseLogger()

	configs.InitialiseConfigs()
	logging.Logger.Info("Configs Loaded")

	err := db.IntialiseDB()
	if err != nil {
		panic(err)
	}
	logging.Logger.Info("DB Initialised")

	err = broker.IntialiseBroker()
	if err != nil {
		panic(err)
	}
	logging.Logger.Info("Broker Intialised")

}

func main() {
	hServer := httpserver.NewHTTPServer(httpserver.ServerConfig{Host: configs.GnCfg.Host, Port: configs.GnCfg.Port})
	logging.Logger.Info("HTTP Server Started", "Host", configs.GnCfg.Host, "Port", configs.GnCfg.Port)
	hServer.ListenAndServe()
}
