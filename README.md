# Go URL Shortener

A simple URL shortener written in Go.

## Features

- Create short URLs
- Custom short names
- Auto-generated short codes
- Redirect using short links
- Prevent duplicate links
- View all links at `/links`

## How It Works

User enters URL  
        ↓  
Server validates URL  
        ↓  
Generate or accept short code  
        ↓  
Store mapping (short → original URL)  
        ↓  
User visits short link  
        ↓  
Server looks up code  
        ↓  
Redirect to original URL