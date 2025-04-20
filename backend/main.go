package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	handler "github.com/markgenuine/tl_hackathon/internal/handlers"
)

const (
	APIURL = "API_URL"
	APIKEY = "API_KEY"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &handler.Config{
		APIURL: os.Getenv(APIURL),
		APIKey: os.Getenv(APIKEY),
	}

	if config.APIURL == "" || config.APIKey == "" {
		log.Fatal("Params for connect openAI not set!")
	}

	http.HandleFunc("/upload", handler.UploadHandler(config))
	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
