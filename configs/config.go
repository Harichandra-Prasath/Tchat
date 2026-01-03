package configs

import "os"

type TableConfig struct {
	UsersTable    string
	SessionsTable string
}

type GeneralConfig struct {

	// Server
	Host string
	Port string

	// RMQ
	RMQUser         string
	RMQPassword     string
	RMQHost         string
	RMQPort         string
	RMQUserExchange string
}

var TbCfg *TableConfig
var GnCfg *GeneralConfig

func InitialiseConfigs() {

	TbCfg = &TableConfig{UsersTable: os.Getenv("USERS_TABLE"), SessionsTable: os.Getenv("SESSIONS_TABLE")}
	GnCfg = &GeneralConfig{
		Host:            os.Getenv("Host"),
		Port:            os.Getenv("PORT"),
		RMQUser:         os.Getenv("RMQ_USER"),
		RMQPassword:     os.Getenv("RMQ_PASSWORD"),
		RMQHost:         os.Getenv("RMQ_HOST"),
		RMQPort:         os.Getenv("RMQ_PORT"),
		RMQUserExchange: os.Getenv("RMQ_USER_EXCHANGE")}
}
