package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getGroupInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")

	// ESTRAZIONE TOKEN
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// LETTURA E VALID CONV
	convid, err := strconv.Atoi(ps.ByName("convid"))
	if err != nil {
		// ID non valido : 400
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// CONTROLLO AUTH UTENTE SU CONVERSAZIONE
	if err := rt.db.IsAuthorized(token, convid); err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	if err := rt.db.IsAGroup(token, convid); err != nil {
		// non è un gruppo : 406 non accettabile
		http.Error(w, "", http.StatusNotAcceptable)
		return
	}

	// RECUPERO INFO GRUPPO
	groupInfo, err := rt.db.GetGroupInfo(convid)
	if err != nil {
		// errore query  400 Bad Request
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(groupInfo)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		rt.baseLogger.WithError(err).Warning("Errore nella scrittura della risposta HTTP")
	}
}
