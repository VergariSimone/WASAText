package api

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Username struct {
	Name string
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	
	// ESTRAZIONE TOKEN
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// LETTURA CONTENUTO BODY
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var u Username
	u.Name = string(body)

	// VALIDAZIONE NOME TRAMITE ESPRESSIONE REGOLARE
	valid, err := regexp.MatchString(`^.+?$`, u.Name)
	if err != nil {

		// ESPRESSIONE REGOLARE NON VALIDA
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	if !valid {

		// PATTERN NON RISPETTATO
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// SET USERNAME NEL DATABASE
	err = rt.db.SetMyUsername(u.Name, token)
	if err != nil {
		if err.Error() == "name already used" {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// OK
	w.WriteHeader(http.StatusNoContent)
}
