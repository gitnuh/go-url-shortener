package main

import (
	"fmt"
	"net/http"
	"strings"
)

// homepage -> shows form
func homeHandler(w http.ResponseWriter, r *http.Request) {

	// if path not "/" try redirect
	if r.URL.Path != "/" {
		redirectHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	// check if we got short code
	short := r.URL.Query().Get("short")

	// check error message
	errorMsg := r.URL.Query().Get("error")

	fmt.Fprintln(w, `<h1>URL Shortener</h1>	`)
	fmt.Fprintln(w, "<a href='/links'>View All Links</a><br><br>")

	// show error if exists
	if errorMsg != "" {
		fmt.Fprintln(w, "<p style='color:red;'>"+errorMsg+"</p>")
	}

	fmt.Fprintln(w, `
	<form action="/create" method="POST">

	Enter URL:<br>
	<input type="text" name="url" size="50"><br><br>

	Short Name (optional):<br>
	<input type="text" name="short"><br><br>

	<button type="submit">Create</button>

	</form>
	`)

	// show result if link created
	if short != "" {

		fmt.Fprintln(w, "<hr>")
		fmt.Fprintln(w, "<h3>short url created</h3>")

		link := "http://" + r.Host + "/" + short
		fmt.Fprintln(w, "<a href='"+link+"'>"+link+"</a>")
	}
}


// create new short url
func createHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		fmt.Fprintln(w, "invalid request")
		return
	}

	originalURL := r.FormValue("url")
	shortName := r.FormValue("short")

	// basic url check- https or http
	if !strings.HasPrefix(originalURL, "http") {
		http.Redirect(w, r, "/?error=invalid+url+(add+http+or+https)", http.StatusSeeOther)
		return
	}

	// check if url actually exists irl
	if !urlReachable(originalURL) {
		http.Redirect(w, r, "/?error=link+does+not+exist", http.StatusSeeOther)
		return
	}

	// check if this url already has a code
	for s, url := range urlMap {

		if url == originalURL {

			if shortName != "" && shortName != s {

				w.Header().Set("Content-Type", "text/html")

				link := "http://" + r.Host + "/" + s
				fmt.Fprintln(w, "<h3>code already exists for this url</h3>")
				fmt.Fprintln(w, "<a href='"+link+"'>"+link+"</a>")
				return
			}

			// reuse existing short code
			http.Redirect(w, r, "/?short="+s, http.StatusSeeOther)
			return
		}
	}

	// generate if empty
	if shortName == "" {
		for {
			shortName = generateShortCode()
			if _, exists := urlMap[shortName]; !exists {
				break
			}
		}
	}

	// check duplicate short name
	if _, exists := urlMap[shortName]; exists {
		http.Redirect(w, r, "/?error=short+name+already+exists", http.StatusSeeOther)
		return
	}

	// store
	urlMap[shortName] = originalURL

	http.Redirect(w, r, "/?short="+shortName, http.StatusSeeOther)
}

// redirect to original
func redirectHandler(w http.ResponseWriter, r *http.Request) {

	// remove "/" from path
	shortName := strings.TrimPrefix(r.URL.Path, "/")

	originalURL, exists := urlMap[shortName]

	if !exists {
		fmt.Fprintln(w, "short url not found")
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}


// show all existing links
func linksHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintln(w, "<h1>Existing Short Links</h1>")
	fmt.Fprintln(w, "<hr>")

	if len(urlMap) == 0 {
		fmt.Fprintln(w, "<p>No links created yet.</p>")
		return
	}

	fmt.Fprintln(w, "<table border='1' cellpadding='5'>")
	fmt.Fprintln(w, "<tr><th>Short Code</th><th>Original URL</th></tr>")

	for short, url := range urlMap {

		shortLink := "http://" + r.Host + "/" + short

		fmt.Fprintln(w, "<tr>")
		fmt.Fprintln(w, "<td><a href='"+shortLink+"'>"+short+"</a></td>")
		fmt.Fprintln(w, "<td>"+url+"</td>")
		fmt.Fprintln(w, "</tr>")
	}

	fmt.Fprintln(w, "</table>")
}