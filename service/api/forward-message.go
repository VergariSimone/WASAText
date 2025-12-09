package api

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// ESTRAZIONE TOKEN
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// PARAM URL
	convid, err := strconv.Atoi(ps.ByName("convid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mesid, err := strconv.Atoi(ps.ByName("mesid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// LETTURA BODY (CORPO)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(body) < 1 {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	target := string(body)
	var conidForw int

	// CONTROLLO SE E' UN NUMERO
	isNumeric := regexp.MustCompile(`^\d+$`).MatchString(target)
	if isNumeric {
		conidForw, err = strconv.Atoi(target)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// AUTORIZZAZIONE
		if err := rt.db.IsAuthorized(token, conidForw); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	} else {
		// CREA O RECUPERA CONVERSAZIONE TRAMITE USERNAME
		conidForw, err = rt.db.SetConversation(token, "conversation", target, "", nil)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// AUTORIZZAZIONE SORGENTE E DESTINATARIO
	if err := rt.db.IsAuthorized(token, convid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err := rt.db.IsAuthorized(token, conidForw); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// INOLTRO EFFETTIVO
	if err := rt.db.ForwardMessage(token, convid, mesid, conidForw); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := struct {
		ConID int `json:"conID"`
	}{
		ConID: conidForw,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonData); err != nil {
		rt.baseLogger.WithError(err).Warning("Errore durante la scrittura della risposta HTTP")
	}
}
