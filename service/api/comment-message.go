package api

import(
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	
	// HEADER JSON 
	w.Header().Set("content-type", "application/json")

	// OK 
	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode("") // SCRIVE STRINGA VUOTA COME JSON
}