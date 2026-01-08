package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// ESTRAZIONE TOKEM
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// PARSING ID CONV
	conID, err := parseConversationID(ps.ByName("convid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // ID NON VALIDO
		return
	}

	if !rt.isGroup(token, conID) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// USCITA GRUPPO
	if !rt.leaveGroupDB(token, conID) {
		w.WriteHeader(http.StatusBadRequest) // ERRORE USCITA
		return
	}

	// OK
	w.WriteHeader(http.StatusNoContent)
}

// FUNZIONI DI SUPORTO

// CONVERSIONE STRING IN INT
func parseConversationID(param string) (int, error) {
	return strconv.Atoi(param)
}

// CONTROLLA SE CONVERSAZIONE È UN GRUPPO
func (rt *_router) isGroup(token string, conID int) bool {
	err := rt.db.IsAGroup(token, conID)
	return err == nil
}

// FUNZIONE CHE EFFETTUA L'USCITA DAL GRUPPO
func (rt *_router) leaveGroupDB(token string, conID int) bool {
	err := rt.db.LeaveGroup(token, conID)
	return err == nil
}
