package application

import (
	"block/domain"
	"encoding/json"
	"net/http"
)

type AddRequest struct {
	Value int64 `json:"value"`
}

type AddResponse struct {
	Ack bool `json:"ack"`
}

func (app *Config) AddHandler(w http.ResponseWriter, r *http.Request) {
	var addRequest AddRequest
	if err := json.NewDecoder(r.Body).Decode(&addRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	job := &domain.AddJob{
		State:       app.State,
		Value:       addRequest.Value,
		RedisClient: app.RedisClient,
	}

	app.WorkerPool.Submit(job)

	bytes, err := json.Marshal(AddResponse{Ack: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(bytes)
}

type GetRequest struct {
	Format string `json:"format"`
}

type GetResponse struct {
	Res any `json:"res"`
}

func (app *Config) GetHandler(w http.ResponseWriter, r *http.Request) {
	var getRequest GetRequest
	if err := json.NewDecoder(r.Body).Decode(&getRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := app.State.Get(getRequest.Format)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bytes, err := json.Marshal(GetResponse{Res: result})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(bytes)
}

type TargetRequest struct {
	Target int64 `json:"target"`
}

type TargetResponse struct {
	Res []int64 `json:"res"`
}

func (app *Config) TargetHandler(w http.ResponseWriter, r *http.Request) {
	var targetRequest TargetRequest
	if err := json.NewDecoder(r.Body).Decode(&targetRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := app.State.Find(targetRequest.Target)

	bytes, err := json.Marshal(TargetResponse{Res: result})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(bytes)
}
