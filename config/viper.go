package config

import (
	"log"

	"github.com/spf13/viper"
)

type EnvVars struct {
	PORT        string `mapstructure:"PORT"`
	MONGODBURI  string `mapstructure:"MONGODB_URI"`
	MONGODBNAME string `mapstructure:"MONGODB_NAME"`
	POSTGRESURI string `mapstructure:"POSTGRES_URI"`
	REDISHOST   string `mapstructure:"REDIS_HOST"`
	REDISDB     int    `mapstructure:"REDIS_DB"`
}

func LoadConfig() (config EnvVars, err error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err)
	}

	return
}
