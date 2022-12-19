package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

type Config struct {
	Server struct {
		Port string
	}
	Smtp struct {
		Host     string
		Port     string
		From     string
		Password string
	}
	Mongo struct {
		Username   string
		Password   string
		Database   string
		Collection string
	}
	Firestore struct {
		ProjectID            string
		SectionsCollectionID string
		UsersCollectionID    string
		CredentialsPath      string
	}
}

func Read() Config {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName(getConfigName())
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
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
	verify(config)
	return config
}

// sample_config is used for deployment to google cloud. It is not used in development.
func getConfigName() string {
	fn := "config.yaml"
	_, err := os.Stat(fn)
	if err != nil {
		log.Println(err)
		log.Println("Using sample_config.yaml")
		return "sample_config.yaml"
	}
	log.Println("Running in development mode")

	return fn
}

func verify(c Config) {
	if c.Server.Port == "" || c.Smtp.Host == "" || c.Smtp.Port == "" || c.Smtp.From == "" || c.Smtp.Password == "" ||
		c.Mongo.Username == "" || c.Mongo.Password == "" || c.Mongo.Database == "" || c.Mongo.Collection == "" ||
		c.Firestore.ProjectID == "" || c.Firestore.SectionsCollectionID == "" || c.Firestore.UsersCollectionID == "" {
		fmt.Println(c.Firestore)
		log.Fatal("required config values are missing")
	}
}
