package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// parse ID Messaggio
	mesID, err := strconv.Atoi(ps.ByName("mesid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// parse Token
	token := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")

	// chiama il db
	count, err := rt.db.UncommentMessage(token, mesID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// controlla il risultato
	if count == 0 {
		// se count è 0, significa che non c'era nessuna reazione da cancellare
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
		return
	}

	// Risposta OK
	w.WriteHeader(http.StatusNoContent)
}
