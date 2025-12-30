package configs

import "os"

type TableConfig struct {
	UsersTable string
}

type GeneralConfig struct {
	Host string
	Port string
}

var TbCfg *TableConfig
var GnCfg *GeneralConfig

func InitialiseConfigs() {

	TbCfg = &TableConfig{UsersTable: os.Getenv("USERS_TABLE")}
	GnCfg = &GeneralConfig{Host: os.Getenv("Host"), Port: os.Getenv("PORT")}
}
