package main

import (
	"log"
	"os"

	"github.com/arekbor/file-manager-server/api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println(err)
		return
	}

	s := api.NewRestApi(os.Getenv("API_ADDR"))
	s.Run()
}
