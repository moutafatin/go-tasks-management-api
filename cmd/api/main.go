package main

import (
	"log"
	"net/http"
)

func main() {
	srv := &http.Server{
		Addr:    ":4000",
		Handler: routes(),
	}

	log.Println("server running on http://localhost:4000")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
