package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/macros-fetcher/handlers"
)

func defaultMsg(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte("Hi, you have made it to the app. I love you, good job."))
}

func getUsdaAPIKey() (string, error) {
	b, err := ioutil.ReadFile("secrets.txt") // just pass the file name
	if err != nil {
		return "", err
	}

	secrets := strings.Split(string(b), "=") // convert content to a 'string'
	fmt.Println(secrets)
	// there is only one secret in the file so secrets[1] is our key
	// todo: make this better for when there are more secrets.
	return secrets[1], nil
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

	UsdaAPIKey, err := getUsdaAPIKey()
	if err != nil {
		log.Fatalf("error getting USDA API key: %v", err)
	}
	ctx := &handlers.Context{
		UsdaAPIKey: UsdaAPIKey,
	}

	http.HandleFunc("/", defaultMsg)
	http.HandleFunc("/v1/ingredients", ctx.IngredientsHandler)
	http.HandleFunc("/v1/macros", ctx.MacrosHandler)

	fmt.Printf("server is listening at %s:%s...\n", host, port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
