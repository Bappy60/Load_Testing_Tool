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
	SECRET_KEY string `mapstructure:"SECRET_KEY"`
	REDIS_HOST string `mapstructure:"REDIS_HOST"`
	REDIS_PORT string `mapstructure:"REDIS_PORT"`
	REDIS_PASS string `mapstructure:"REDIS_PASS"`
}

func InitConfig() *Config {

	viper.AddConfigPath("/app")

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
