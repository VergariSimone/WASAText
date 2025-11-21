package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// ESTRAZIONE TOKEN
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// ESTRAZIONE DEGLI UTENTI
	usernames := r.URL.Query()["users"]
	if len(usernames) < 1 || len(usernames) > 32 {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	// PRENDI ID GRUPPO
	groupIDStr := r.URL.Query().Get("groupID")
	var groupID int
	var err error

	if groupIDStr != "" {
		groupID, err = strconv.Atoi(groupIDStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		// CREA NUOVO GRUPPO
		groupID, err = rt.db.SetGroup(token)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// AGGIUNTA UTENTI AL GRUPPO
	if err := rt.db.AddToGroup(usernames, groupID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// RISPOSTA JSON CON ID DEL GRUPPO
	resp, err := json.Marshal(groupID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		rt.baseLogger.WithError(err).Warning("Errore durante la scrittura della risposta HTTP")
	}
}
