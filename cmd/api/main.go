package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Write([]byte("Hello there"))
	})

	srv := &http.Server{
		Addr:    ":4000",
		Handler: mux,
	}

	log.Println("server running on http://localhost:4000")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
