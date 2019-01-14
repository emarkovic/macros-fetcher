package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/macros-fetcher/handlers"
)

func defaultMsg(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte("Hi, you have made it to the app. I love you, good job."))
}

func main() {

	host := os.Getenv("HOST")
	if len(host) == 0 {
		host = "localhost"
	}

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	http.HandleFunc("/", defaultMsg)
	http.HandleFunc("/v1/ingredients", handlers.IngredientsHandler)

	fmt.Printf("server is listening at %s:%s...\n", host, port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
