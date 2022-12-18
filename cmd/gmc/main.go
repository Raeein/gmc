package main

import (
	"fmt"
	"github.com/Raeein/gmc/config"
	"github.com/Raeein/gmc/server"
	"github.com/Raeein/gmc/webadvisor"
)

func main() {

	cfg := config.Read()
	fmt.Println("Host is:", cfg.Smtp.Host)
	//mongoService := mongodb.New(
	//	cfg.Mongo.Username,
	//	cfg.Mongo.Password,
	//	cfg.Mongo.Database,
	//	cfg.Mongo.Collection,
	//	)
	//defer mongoService.Close()
	//mongoService.Log("Info", "test log entry")

	webadvisorService := webadvisor.New()
	webadvisorService.CheckCourse()
	server.Start(cfg.Server.Port)
	//fmt.Println("Server is running")
	//u := gmc.User{"raeein@aol.com", "Raeein"}
	//email.Send(u, []string{"Java", "Python"}, cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.From, cfg.Smtp.Password)
	//p := &gmc.User{"jason@me.com", "hi"}
}
