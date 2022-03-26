package config

import (
	"gopkg.in/ini.v1"
	"log"
)

type MysqlConfig struct {
	Host     string `ini:"host"`
	Port     string `ini:"port"`
	Dbname   string `ini:"dbname"`
	Username string `ini:"username"`
	Password string `ini:"password"`
}

var Config *MysqlConfig

func init() {
	Config = new(MysqlConfig)
	cfg, err := ini.Load("config.ini")
	if err != nil {
		panic(err)
	}
	if err := cfg.Section("mysql").MapTo(Config); err != nil {
		panic(err)
	}
	log.Println("配置文件加载成功")
}
