package application

import (
	"accumulator/domain"
	"accumulator/infrastructure"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	State      *domain.State
	WorkerPool *infrastructure.WorkerPool
	Limiter    *RateLimiter
}

// Adds the routes to the application
func (app *Config) Routes() http.Handler {
	mux := mux.NewRouter()

	mux.Handle("/add", app.Limiter.RateLimitMiddleware(http.HandlerFunc(app.AddHandler))).Methods("POST")
	mux.HandleFunc("/state", app.GetHandler).Methods("GET")
	mux.HandleFunc("/tsp", app.TargetHandler).Methods("GET")

	return mux
}
