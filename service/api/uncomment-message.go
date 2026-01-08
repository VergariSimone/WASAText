package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Parse ID Messaggio
	// Nota: Usiamo "mesid" perché la cancellazione ora è basata sul messaggio
	mesID, err := strconv.Atoi(ps.ByName("mesid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 2. Parse Token
	token := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")

	// 3. Chiama il Database
	// Passiamo (token, mesid) come richiesto dalla tua nuova firma DB
	count, err := rt.db.UncommentMessage(token, mesID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 4. Controlla il risultato
	if count == 0 {
		// Se count è 0, significa che non c'era nessuna reazione da cancellare
		// (l'utente non aveva messo like o il messaggio non esiste)
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
		return
	}

	// 5. Risposta OK (204 No Content)
	w.WriteHeader(http.StatusNoContent)
}
