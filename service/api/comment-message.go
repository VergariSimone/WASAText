package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// HEADER JSON
	w.Header().Set("content-type", "application/json")

	// OK
	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode("") // SCRIVE STRINGA VUOTA COME JSON
}
