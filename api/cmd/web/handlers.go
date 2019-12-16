package main

import (
	"fmt"
	"html"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Hello there %q", html.EscapeString(r.URL.Path))))
}

func (app *application) products(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Products %q", html.EscapeString(r.URL.Path))))
}
