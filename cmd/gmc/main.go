package main

import (
	"context"
	"github.com/Raeein/gmc/config"
	"github.com/Raeein/gmc/mongodb"
	"github.com/Raeein/gmc/server"
	"github.com/Raeein/gmc/webadvisor"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		log.Println("received interrupt, shutting down")
		cancel()
	}()

	cfg := config.Read()
	mongoService := mongodb.New(
		cfg.Mongo.Username,
		cfg.Mongo.Password,
		cfg.Mongo.Database,
		cfg.Mongo.Collection,
	)
	//defer mongoService.Close()
	//mongoService.Log("Info", "test log entry")

	webadvisorService := webadvisor.New()
	//firestoreService, err := firestore.New(
	//	ctx,
	//	cfg.Firestore.ProjectID,
	//	cfg.Firestore.SectionsCollectionID,
	//	cfg.Firestore.UsersCollectionID,
	//	cfg.Firestore.CredentialsPath,
	//)
	s := server.New(cfg.Server.Port, webadvisorService, mongoService)
	s.Start(ctx, cfg.Server.Port)
	//fmt.Println("Service is running")
	//u := gmc.User{"raeein@aol.com", "Raeein"}
	//email.Send(u, []string{"Java", "Python"}, cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.From, cfg.Smtp.Password)
	//p := &gmc.User{"jason@me.com", "hi"}
}
