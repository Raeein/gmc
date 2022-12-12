package main

import (
	"github.com/Raeein/gmc"
	"github.com/Raeein/gmc/config"
	"github.com/Raeein/gmc/email"
)

func main() {
	cfg := config.Read()
	u := gmc.User{"test@aol.com", "Raeein"}
	email.Send(u, cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.From, cfg.Smtp.Password)
	//p := &gmc.User{"jason@me.com", "hi"}

}
