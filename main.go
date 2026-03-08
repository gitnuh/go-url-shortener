package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// map storing short -> original
var urlMap = make(map[string]string)

func main() {

	// seed random once
	rand.Seed(time.Now().UnixNano())

	// routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/links", linksHandler)

	fmt.Println("server running at http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}