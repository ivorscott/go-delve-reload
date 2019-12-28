package mid

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/ivorscott/go-delve-reload/internal/platform/web"
	"github.com/pkg/errors"
)

// Panics recovers from panics and converts the panic to an error so it is
// reported in Metrics and handled in Errors.
func Panics(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(after web.Handler) web.Handler {

		h := func(w http.ResponseWriter, r *http.Request) (err error) {

			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if r := recover(); r != nil {
					err = errors.Errorf("panic: %v", r)

					// Log the Go stack trace for this panic'd goroutine.
					log.Printf("%s", debug.Stack())
				}
			}()

			// Call the next Handler and set its return value in the err variable.
			return after(w, r)
		}

		return h
	}

	return f
}
