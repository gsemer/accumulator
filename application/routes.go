package application

import (
	"block/domain"
	"block/infrastructure"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	State      *domain.State
	WorkerPool *infrastructure.WorkerPool
}

// Adds the routes to the application
func (app *Config) Routes() http.Handler {
	mux := mux.NewRouter()

	mux.HandleFunc("/add", app.AddHandler).Methods("POST")
	mux.HandleFunc("/state", app.GetHandler).Methods("GET")
	mux.HandleFunc("/tsp", app.TargetHandler).Methods("GET")

	return mux
}
