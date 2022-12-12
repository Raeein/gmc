package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Smtp struct {
		Host     string
		Port     int
		From     string
		Password string
	}
}

func Read() Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal(err)
		} else {
			log.Fatal(err)
		}
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	return config
}
