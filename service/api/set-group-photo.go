package api

import (
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// ESTRAZIONE TOKEN 
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// PARSING CONV ID
	convid, err := strconv.Atoi(ps.ByName("convid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// VERIFICA DEL TIPO
	if ct := r.Header.Get("Content-type"); ct != "image/gif" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// LETTURA FILE
	fileBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// VALIDAZIONE GIF
	isGifHeader, err := regexp.Match(`^GIF`, fileBytes)
	if err != nil || !isGifHeader {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	// RIDUCIAMO LUNGHEZZA CONTROLLANDO IL BUFFER
	fileLen := len(fileBytes)
	if fileLen < 4 || fileLen > 524288 {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	// VERIFICA CONVERSIONE SIA GRUPPO
	if err := rt.db.IsAGroup(token, convid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// VERIFICA AUTH
	if err := rt.db.IsAuthorized(token, convid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// SALVATAGGIO
	if err := rt.db.SetGroupPhoto(convid, fileBytes); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 204 OK
	w.WriteHeader(http.StatusNoContent)
}
