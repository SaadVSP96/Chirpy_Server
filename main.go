package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux() // router
	mux.Handle("/", http.FileServer(http.Dir(".")))
	server := &http.Server{
		Addr:    ":8080", // Bind to localhost:8080
		Handler: mux,     // Use our empty mux
	}

	log.Println("server running on :8080")
	log.Fatal(server.ListenAndServe())

}
