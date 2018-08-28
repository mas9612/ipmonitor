package main

import (
	"log"
	"net/http"

	"github.com/mas9612/ipmonitor"
)

func main() {
	err := ipmonitor.OpenDB("test.db")
	if err != nil {
		log.Fatalln("Failed to open DB.")
	}
	ipmonitor.InitDB()
	defer ipmonitor.Conn.DB.Close()

	handler := ipmonitor.NewHTTPHandler()
	log.Fatalln(http.ListenAndServe("localhost:8080", handler))
}
