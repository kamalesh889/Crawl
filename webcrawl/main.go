package main

import (
	"fmt"
	"net/http"

	"main.go/Rout"
)

func main() {
	fmt.Println("Api started")
	mux := Rout.Router()
	http.Handle("/", mux)
	http.ListenAndServe(":8080", mux)
}
