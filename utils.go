package main

import (
	"math/rand"
	"net/http"
	"time"
)

// gen random short code
func generateShortCode() string {

	chars := "abcdefghijklmnopqrstuvwxyz0123456789"

	code := ""

	for i := 0; i < 5; i++ {
		code += string(chars[rand.Intn(len(chars))])
	}

	return code
}

// check if site actually responds
func urlReachable(link string) bool {

	client := http.Client{
		Timeout: 3 * time.Second, 
	}

	resp, err := client.Head(link)

	if err != nil {
		return false
	}

	defer resp.Body.Close()
	return resp.StatusCode < 400
}