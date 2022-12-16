package server

import (
	"fmt"
	"log"
	"net/http"
)

func Start() {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/register", register)
	fmt.Println("Visit localhost:8080")
	http.ListenAndServe("localhost:8080", nil)
}

func ping(rw http.ResponseWriter, r *http.Request) {
	_, err := rw.Write([]byte("<h1>Im Alive</h1>"))
	if err != nil {
		log.Println(err)
	}
}

func register(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case http.MethodPost:
		fmt.Println("It was a post")
		rw.Write([]byte("<h1>Registered</h1>"))
		fmt.Println(r.Body)
	case http.MethodOptions:
		rw.Header().Set("Allow", "GET, POST, OPTIONS")
		rw.WriteHeader(http.StatusNoContent)

	default:
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func get() {
	resp, err := http.Get("https://aol.com")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp)
}
