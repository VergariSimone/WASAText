package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// ESTRAZIONE DEL TOKEN DALL HEADER
	token := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")

	// CONVERSIONE DA STRINGA A INTERO DI ID MESSAGGIO
	messageID, err := strconv.Atoi(ps.ByName("mesid"))
	if err != nil {
		// SE NON VALIDO 400
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// RICHIAMA MET DATABASE PER ELIMINARE COMMENTO DAL MESSAGGIO
	_, err = rt.db.UncommentMessage(token, messageID)
	if err != nil {

		// SE CE ERRORE 400
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 204 OK
	w.WriteHeader(http.StatusNoContent)

}
