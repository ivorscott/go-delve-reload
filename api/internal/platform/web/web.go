package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/go-chi/chi"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// Values carries information about each request.
type Values struct {
	StatusCode int
	Start      time.Time
}

type Handler func(http.ResponseWriter, *http.Request) error

// App is the entry point for all applications
type App struct {
	mux      *chi.Mux
	log      *log.Logger
	mw       []Middleware
	shutdown chan os.Signal
}

// New App constructs internal state for a new  app
func NewApp(shutdown chan os.Signal, logger *log.Logger, mw ...Middleware) *App {
	return &App{
		log:      logger,
		mux:      chi.NewRouter(),
		mw:       mw,
		shutdown: shutdown,
	}
}

// Handle associates a handler function with an HTTP Method and URL pattern.
//
// It converts our custom handler type to the std lib Handler type. It captures
// errors from the handler and serves them to the client in a uniform way.
func (a *App) Handle(method, url string, h Handler) {

	h = wrapMiddleware(a.mw, h)

	fn := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			Start: time.Now(),
		}

		ctx := r.Context()                          // get original context
		ctx = context.WithValue(ctx, KeyValues, &v) // create a new context with new key/value
		// you can't directly update a request context
		r = r.WithContext(ctx) // create a new request and pass context

		// Call the handler and catch any propagated error.
		if err := h(w, r); err != nil {
			a.log.Printf("ERROR : unhandled error\n %+v", err)
			if IsShutdown(err) {
				a.SignalShutdown()
			}
		}
	}

	a.mux.MethodFunc(method, url, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	a.mux.ServeHTTP(w, r)
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.log.Println("error returned from handler indicated integrity issue, shutting down service")
	a.shutdown <- syscall.SIGSTOP
}
