package main

import (
	"fmt"
	"github.com/Raeein/gmc/config"
)

func main() {
	cfg := config.Read()
	fmt.Println(cfg)
	//mc := mongodb.New(
	//	cfg.Mongo.Username,
	//	cfg.Mongo.Password,
	//	cfg.Mongo.Database,
	//	cfg.Mongo.Collection,
	//)
	//mc.Log("Info", "test log entry")
	//defer mc.Close()
	//server.Start()
	//fmt.Println("Server is running")
	//u := gmc.User{"raeein@aol.com", "Raeein"}
	//email.Send(u, []string{"Java", "Python"}, cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.From, cfg.Smtp.Password)
	//p := &gmc.User{"jason@me.com", "hi"}
}
