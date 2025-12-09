package api

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// ESTRAZIONE TOKEN
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// PARSING CONV ID
	convid, err := strconv.Atoi(ps.ByName("convid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// RECUPERO NOME
	name := r.URL.Query().Get("name")

	// VALIDAZIONE DEL PATTERN
	matches, err := regexp.MatchString(`.*`, name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !matches {
		// NON RISPETTATO PATTERN
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	// VALIDITA' LUNGH
	if len(name) < 1 || len(name) > 32 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// VERIFICA SE CONV IS GRUPPO
	if err := rt.db.IsAGroup(token, convid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// VERIFICA AUTH UTENTE
	if err := rt.db.IsAuthorized(token, convid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// NUOVO NOME DEL GRUPPO
	if err := rt.db.SetGroupName(convid, name); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/// OK 204
	w.WriteHeader(http.StatusNoContent)
}
