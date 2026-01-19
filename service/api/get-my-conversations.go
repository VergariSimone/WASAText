package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)
// restituisce la lista delle chat ordinata per timestamp decrescente, includendo l'anteprima dell'ultimo messaggio e il contatore dei messaggi non letti (toread)
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")

	// ESTRAZIONE TOKEN
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// RECUPERO CONV UTENTE
	convs, err := rt.db.GetMyConversations(token)
	if err != nil {

		http.Error(w, "", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(convs)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		rt.baseLogger.WithError(err).Warning("Errore nella scrittura della risposta HTTP")
	}
}
