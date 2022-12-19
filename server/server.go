package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/Raeein/gmc/mongodb"
	"github.com/Raeein/gmc/webadvisor"
	"log"
	"net/http"
	"text/template"
	"time"
)

type Service struct {
	server            *http.Server
	webadvisorService webadvisor.WebAdvisor
	mongoService      mongodb.Logger
}

func New(port string, wa webadvisor.WebAdvisor, ms mongodb.Logger) Service {
	localPort := fmt.Sprintf(":%s", port)
	return Service{
		server:            &http.Server{Addr: localPort},
		webadvisorService: wa,
		mongoService:      ms,
	}
}

func (s Service) Start(ctx context.Context, port string) {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", home)
	http.HandleFunc("/home", home)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/register", s.register)

	errChan := make(chan error)
	go func() {
		log.Printf("Visit localhost:%s", port)
		err := s.server.ListenAndServe()
		errChan <- err
	}()

	select {
	case err := <-errChan:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.server.Shutdown(ctx); err != nil {
			log.Fatal("server shutdown error: %w", err)
		}
		log.Println("server shutdown complete")
	}
}

func home(rw http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("static/templates/index.gohtml"))
	err := tpl.Execute(rw, nil)
	if err != nil {
		log.Println(err)
	}
}

func ping(rw http.ResponseWriter, r *http.Request) {
	_, err := rw.Write([]byte("<h1>Im Alive</h1>"))
	if err != nil {
		log.Println(err)
	}
}

func (s Service) register(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case http.MethodPost:
		fmt.Println("It was a post")
		err := validateRegistration(&s.webadvisorService, r.Body)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println("Found the course")
		rw.WriteHeader(http.StatusOK)
	case http.MethodOptions:
		rw.Header().Set("Allow", "GET, POST, OPTIONS")
		rw.WriteHeader(http.StatusNoContent)

	default:
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
	}
}
