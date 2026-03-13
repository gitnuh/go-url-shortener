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

	fmt.Fprintln(w, `
	<html>
	<head>
	<title>URL Shortener</title>
	<style>

	body{
		margin:0;
		height:100vh;
		display:flex;
		justify-content:center;
		align-items:center;
		background:linear-gradient(135deg,#3a7bd5,#8e44ad);
		font-family:Arial, Helvetica, sans-serif;
		color:white;
	}

	.container{
	background:rgba(255,255,255,0.9);
	padding:60px 80px;
	border-radius:18px;
	text-align:center;
	box-shadow:0 15px 40px rgba(0,0,0,0.25);
	width:500px;
	max-width:90%;
	color:#222;
}
	h1{
		font-size:40px;
		font-weight:bold;
		margin-bottom:20px;
	}

	input{
	padding:12px;
	border-radius:8px;
	border:2px solid #5b6ee1;
	width:420px;
	font-size:16px;
	margin-top:8px;
	background:#f4f6ff;
	color:#333;
}

	button{
	padding:12px 28px;
	border:none;
	border-radius:8px;
	background:#5b6ee1;
	color:white;
	font-weight:bold;
	cursor:pointer;
	font-size:16px;
}

	button:hover{
		background:#8e44ad;
	}

	a{
	color:#4b63d1;
	font-weight:bold;
	text-decoration:none;
}

	.error{
		color:d63384;
		font-weight:bold;
	}
	.label{
	text-align:left;
	width:420px;
	margin:auto;
	font-weight:bold;
	font-weight:500;
	margin-top:10px;
}

	.success{
	color:black;
	font-weight:bold;
	font-size:18px;
}

	</style>
	</head>
	<body>

	<div class="container">

	<h1>URL Shortener</h1>

	<a href='/links'>View All Links</a><br><br>
	`)

	// show error if exists
	if errorMsg != "" {
		fmt.Fprintln(w, "<p class='error'>"+errorMsg+"</p>")
	}

	fmt.Fprintln(w, `
	<form action="/create" method="POST">

	<p class="label">Enter URL:</p>
<input type="text" name="url"><br><br>


<p class="label">Short Name (Optional):</p>
<input type="text" name="short"><br><br>

	<button type="submit">Create</button>

	</form>
	`)

	// show result if link created
	if short != "" {

		fmt.Fprintln(w, "<hr>")
		fmt.Fprintln(w, "<h3 class='success'>Short URL Created:</h3>")

		link := "http://" + r.Host + "/" + short
		fmt.Fprintln(w, "<a href='"+link+"'>"+link+"</a>")
	}

	fmt.Fprintln(w, `
	</div>
	</body>
	</html>
	`)
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
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "http://" + originalURL
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

	fmt.Fprintln(w, `
	<html>
	<head>
	<title>Existing Links</title>
	<style>

	body{
		margin:0;
		height:100vh;
		display:flex;
		justify-content:center;
		align-items:center;
		background:linear-gradient(135deg,#3a7bd5,#8e44ad);
		font-family:Arial, Helvetica, sans-serif;
		color:white;
	}

	.container{
		background:rgba(0,0,0,0.35);
		padding:60px 80px;
		border-radius:18px;
		text-align:center;
		box-shadow:0 10px 25px rgba(0,0,0,0.4);
		width:500px;
		max-width:90%;
	}

	table{
		background:white;
		color:black;
		border-collapse:collapse;
		margin-top:20px;
		margin:20px auto;
	}

	th{
		background:#333;
		color:white;
	}

	td,th{
		padding:10px;
		text-align:center;
	}

	a{
		color:#3a7bd5;
		font-weight:bold;
		text-decoration:none;
	}

	</style>
	</head>

	<body>

	<div class="container">

	<h1>Existing Short Links</h1>
	<hr>
	`)

	if len(urlMap) == 0 {
		fmt.Fprintln(w, "<p>No links created yet.</p>")
		fmt.Fprintln(w, "</div></body></html>")
		return
	}

	fmt.Fprintln(w, "<table border='1'>")
	fmt.Fprintln(w, "<tr><th>Short Code</th><th>Original URL</th></tr>")

	for short, url := range urlMap {

		shortLink := "http://" + r.Host + "/" + short

		fmt.Fprintln(w, "<tr>")
		fmt.Fprintln(w, "<td><a href='"+shortLink+"'>"+short+"</a></td>")
		fmt.Fprintln(w, "<td>"+url+"</td>")
		fmt.Fprintln(w, "</tr>")
	}

	fmt.Fprintln(w, "</table>")

	fmt.Fprintln(w, `
	</div>
	</body>
	</html>
	`)
}
