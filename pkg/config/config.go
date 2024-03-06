package config

import (
	"log"

	"github.com/spf13/viper"
)

var GConfig *Config

type Config struct {
	DBUser string `mapstructure:"DBUSER"`
	DBPass string `mapstructure:"DBPASS"`
	DBHost string `mapstructure:"DBHOST"`
	DbName string `mapstructure:"DBNAME"`
	DBPort string `mapstructure:"DBPORT"`
}

func InitConfig() *Config {

	viper.AddConfigPath("/home/vivasoft/bappy/goProjects/Load_Testing_Tool-main")

	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	var config *Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error reading env file", err)
	}

	return config

}

func SetConfig() {
	GConfig = InitConfig()
}
