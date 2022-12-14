package main

import (
	"github.com/Raeein/gmc/config"
	"github.com/Raeein/gmc/mongodb"
)

func main() {
	cfg := config.Read()
	mc := mongodb.New(
		cfg.Mongo.Username,
		cfg.Mongo.Password,
		cfg.Mongo.Database,
		cfg.Mongo.Collection,
	)
	mc.Log("Info", "test log entry")
	mc.Close()
	//u := gmc.User{"raeein@aol.com", "Raeein"}
	//email.Send(u, []string{"Java", "Python"}, cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.From, cfg.Smtp.Password)
	//p := &gmc.User{"jason@me.com", "hi"}
}
