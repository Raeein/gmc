package server

import (
	"fmt"
	"github.com/Raeein/gmc/mongodb"
	"github.com/Raeein/gmc/webadvisor"
	"log"
	"net/http"
	"text/template"
)

type Server struct {
	port              string
	webadvisorService webadvisor.WebAdvisor
	mongoService      mongodb.Logger
}

func New(port string, wa webadvisor.WebAdvisor, ms mongodb.Logger) Server {
	return Server{port: port, webadvisorService: wa, mongoService: ms}
}

func (s Server) Start(port string) {

	fs := http.FileServer(http.Dir("static"))
	http.HandleFunc("/", home)
	http.HandleFunc("/home", home)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/register", s.register)
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Printf("Visit localhost:%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
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

func (s Server) register(rw http.ResponseWriter, r *http.Request) {
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
