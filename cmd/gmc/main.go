package main

import (
	"github.com/Raeein/gmc/config"
	"github.com/Raeein/gmc/mongodb"
)

func main() {
	cfg := config.Read()
	mc := mongodb.New(cfg.Mongo.Username, cfg.Mongo.Password)
	mc.Log(cfg.Mongo.Database, cfg.Mongo.Collection)
	//u := gmc.User{"raeein@aol.com", "Raeein"}
	//email.Send(u, []string{"Java", "Python"}, cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.From, cfg.Smtp.Password)
	//p := &gmc.User{"jason@me.com", "hi"}
}
