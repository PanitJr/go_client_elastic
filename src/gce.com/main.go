package main

import (
	"log"
	"net/http"

	"go_client_elastic/src/gce.com/sub"
)


func main() {

	router := sub.NewRouter()

	log.Fatal(http.ListenAndServe(":8082", router))
}
