package main

import (
	"flag"
	"log"

	"github.com/arekbor/file-manager-server/api"
	"github.com/joho/godotenv"
)

var (
	listenAddr = flag.String("listenAddr", ":8080", "listen addr of rest api")
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println(err)
		return
	}

	s := api.NewRestApi(*listenAddr)
	s.Run()
}
