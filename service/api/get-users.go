package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")

	// ESTRAZIONE TOKEN
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// LETTURA PARAM DELLA QUERY
	fragment := r.URL.Query().Get("search")
	showgroups := r.URL.Query().Get("showgroups")

	groups := false

	if showgroups != "" && showgroups != "false" && showgroups != "true" {
		// Il parametro ha un valore non accettabile : 406
		http.Error(w, "", http.StatusNotAcceptable)
		return
	}
	if showgroups == "true" {
		groups = true
	}

	// QUERY AL DB
	users, err := rt.db.GetUsers(token, fragment, groups)
	if err != nil {
		// Nessun risultato :204 No Content
		if err.Error() == "sql: no rows in result set" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// Errori generici : 400 Bad Request
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		rt.baseLogger.WithError(err).Warning("Errore nella scrittura della risposta HTTP")
	}
}
