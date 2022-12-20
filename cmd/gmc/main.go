package main

import (
	"context"
	"github.com/Raeein/gmc/config"
	"github.com/Raeein/gmc/email"
	"github.com/Raeein/gmc/firestore"
	"github.com/Raeein/gmc/mongodb"
	"github.com/Raeein/gmc/poll"
	"github.com/Raeein/gmc/server"
	"github.com/Raeein/gmc/webadvisor"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	ticker := time.NewTicker(30 * time.Minute)

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

	firestoreService, err := firestore.New(
		ctx,
		cfg.Firestore.ProjectID,
		cfg.Firestore.SectionsCollectionID,
		cfg.Firestore.UsersCollectionID,
		cfg.Firestore.CredentialsPath,
	)
	if err != nil {
		log.Fatalf("failed to create firestore service: %v", err)
	}
	s := server.New(ctx, webadvisorService, firestoreService, mongoService, cfg.Server.Port)
	if err != nil {
		log.Fatalf("failed to convert port string to int: %v", err)
	}
	emailService := email.New(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.From, cfg.Smtp.Password)
	trigger := poll.New(webadvisorService, firestoreService, emailService, mongoService)

	go func() {
		for {
			select {
			case <-ticker.C:
				go trigger.Trigger(ctx)
			case <-ctx.Done():
				ticker.Stop()
				cancel()
			case <-sig:
				log.Println("received interrupt, shutting down")
				ticker.Stop()
				cancel()
			}
		}
	}()

	s.Start(ctx, cfg.Server.Port)
	//u := gmc.User{"raeein@aol.com", "Raeein"}
	//email.Send(u, []string{"Java", "Python"}, cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.From, cfg.Smtp.Password)
	//p := &gmc.User{"jason@me.com", "hi"}
}
