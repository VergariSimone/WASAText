package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// ESTRAZIONE TOKEN
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// LETTURA E VALIDAZIONE ID CONV
	convid, err := strconv.Atoi(ps.ByName("convid"))
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// LETTURA PARAMETRO START
	startInt, err := parseQueryInt(r, "start", 0)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if startInt < 0 {
		http.Error(w, "", http.StatusNotAcceptable)
		return
	}

	// lETTURA PARAM LIMIT
	limitInt, err := parseQueryInt(r, "limit", 50)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if limitInt < 50 || limitInt > 100 {
		http.Error(w, "", http.StatusNotAcceptable)
		return
	}

	// CONTROLLO AUTH
	if err := rt.db.IsAuthorized(token, convid); err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	// RECUPEROP DA DB
	convs, err := rt.db.GetConversation(convid, startInt, limitInt, token)
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

// Funzione di supporto per leggere un parametro numerico dalla query string
// Restituisce un valore di default se il parametro non esiste
func parseQueryInt(r *http.Request, key string, defaultVal int) (int, error) {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal, nil
	}
	return strconv.Atoi(val)
}
