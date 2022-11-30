package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	go startHTTP()

	// Block forever
	var wg sync.WaitGroup

	wg.Add(1)
	wg.Wait()
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "www/swagger.json")
}

func startHTTP() {
	// Serve the swagger,
	mux := http.NewServeMux()
	mux.HandleFunc("/swagger.json", serveSwagger)

	fs := http.FileServer(http.Dir("www/swagger-ui"))
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", fs))

	log.Println("REST server ready...")
	log.Println("Serving Swagger at: http://localhost:8080/swagger-ui/")

	const timeout = 5

	srv := &http.Server{ //nolint:exhaustruct
		Addr:         "localhost:8000",
		Handler:      mux,
		ReadTimeout:  timeout * time.Second,
		WriteTimeout: timeout * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
