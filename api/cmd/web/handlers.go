package main

import (
	"fmt"
	"html"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(fmt.Sprintf("Hello there.. %q", html.EscapeString(r.URL.Path))))
	if err != nil {
		app.errorLog.Fatal(err)
	}
}

func (app *application) products(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(fmt.Sprintf("Products %q", html.EscapeString(r.URL.Path))))
	if err != nil {
		app.errorLog.Fatal(err)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		println("ping failed")
	}
}
