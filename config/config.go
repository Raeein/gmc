package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Smtp struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		From     string `yaml:"from"`
		Password string `yaml:"password"`
	} `yaml:"Smtp"`
	Mongo struct {
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		Database   string `yaml:"database"`
		Collection string `yaml:"collection"`
	} `yaml:"Mongo"`
}

func Read() Config {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

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
	fmt.Println(config)
	return config
}
