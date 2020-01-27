package api

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
)

// home displays a welcome message
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(fmt.Sprintf("Hello there.. %q", html.EscapeString(r.URL.Path))))
	if err != nil {
		app.errorLog.Fatal(err)
	}
}

// showProducts displays products from the database
func (app *application) showProducts(w http.ResponseWriter, r *http.Request) {
	products := app.products.GetAll()

	js, err := json.Marshal(products)
	if err != nil {
		app.errorLog.Printf("%s %d", err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(js)
	if err != nil {
		app.errorLog.Fatal(err)
	}
}
