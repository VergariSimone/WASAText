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

	//  VERIFICA DEL TIPO (GIF, PNG o JPEG)
	ct := r.Header.Get("Content-type")
	if ct != "image/gif" && ct != "image/png" && ct != "image/jpeg" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// LETTURA FILE
	fileBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// VALIDAZIONE FORMATO
	isGifHeader, _ := regexp.Match(`^GIF`, fileBytes)
	isPngHeader, _ := regexp.Match(`^\x89PNG`, fileBytes)
	isJpgHeader, _ := regexp.Match(`^\xff\xd8`, fileBytes)

	if (!isGifHeader && !isPngHeader && !isJpgHeader) {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	// RIDUCIAMO LUNGHEZZA CONTROLLANDO IL BUFFER
	fileLen := len(fileBytes)
	if fileLen < 4 || fileLen > 1048576 {
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