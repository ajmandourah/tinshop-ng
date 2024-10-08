package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ajmandourah/tinshop/repository"
)

type endpoint struct {
}

// New returns a new api
func New() repository.API {
	return &endpoint{}
}

func (e *endpoint) Stats(w http.ResponseWriter, stats repository.StatsSummary) {
	jsonResponse, jsonError := json.Marshal(stats)

	if jsonError != nil {
		log.Println("[API] Unable to encode JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}
