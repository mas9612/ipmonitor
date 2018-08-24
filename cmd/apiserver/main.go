package main

import (
	"log"
	"net/http"

	"github.com/mas9612/ipmonitor"
)

func main() {
	ipmonitor.InitDB()

	handler := ipmonitor.NewHTTPHandler()
	log.Fatalln(http.ListenAndServe("localhost:8080", handler))
}
