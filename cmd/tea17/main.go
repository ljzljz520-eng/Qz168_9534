package main

import (
	"flag"
	"log"
	"net/http"
	"tea17"
)

func main() {
	address := flag.String("listen", ":8080", "HTTP listen address")
	database := flag.String("db", "./data/tea17.db", "bbolt database path")
	flag.Parse()
	app, err := tea17.Open(*database)
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()
	log.Printf("tea17 listening on %s", *address)
	if err = http.ListenAndServe(*address, app.API.Handler()); err != nil {
		log.Fatal(err)
	}
}
