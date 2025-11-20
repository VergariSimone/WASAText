package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// HEADER JSON
	w.Header().Set("Content-Type", "application/json")

	// OK 200
	w.WriteHeader(200)

	// RISPOSTA
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "messaggio eliminato",
	})
}
