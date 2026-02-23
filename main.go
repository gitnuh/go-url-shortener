package main

import (
	"fmt"
	"net/http"
)

// map to store short name -> original URL
var urlMap = make(map[string]string)

// handler function to redirect
func redirectHandler(w http.ResponseWriter, r *http.Request) {

	// remove "/" from path
	shortName := r.URL.Path[1:]

	// find original URL
	originalURL, exists := urlMap[shortName]
	if !exists {
		fmt.Fprintln(w, "Short URL not found")
		return
	}

	// redirect to original URL
	http.Redirect(w, r, originalURL, http.StatusFound)
}

func main() {

	var originalURL string
	var shortName string

	// take input from terminal
	fmt.Print("Enter original URL: ")
	fmt.Scanln(&originalURL)

	fmt.Print("Enter preferred short name: ")
	fmt.Scanln(&shortName)

	// store in map
	urlMap[shortName] = originalURL

	// start server
	http.HandleFunc("/", redirectHandler)

	fmt.Println("\nShort URL created successfully!")
	fmt.Println("Click here:", "http://localhost:8080/"+shortName)
	fmt.Println("Server is running...")

	// run server
	http.ListenAndServe(":8080", nil)
}