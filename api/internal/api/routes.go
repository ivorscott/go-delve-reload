package api

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"github.com/rs/cors"
)

// routes connects the api endpoints, middleware the corresponding handlers
func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, app.secureHeaders)
	dynamicMiddleware := alice.New(app.authenticate)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://client.local"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization"},
	})

	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/products", dynamicMiddleware.ThenFunc(app.showProducts))
	return standardMiddleware.Then(c.Handler(mux))
}
