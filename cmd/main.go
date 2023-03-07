package main

import (
	"flag"

	"github.com/arekbor/file-manager-server/api"
	"github.com/joho/godotenv"
)

var (
	listenAddr = flag.String("listenAddr", ":8000", "listen addr of rest api")
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		panic(err)
	}

	s := api.NewRestApi(*listenAddr)
	s.Run()
}
