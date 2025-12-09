package api

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Userphoto struct {
	Photo []byte
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// ESTRAZIONE TOKEN
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// CONTROLLO DEL CONTENT-TYPE
	if r.Header.Get("Content-Type") != "image/gif" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// LETTURA FILE IMMAGINE
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var u Userphoto
	u.Photo = data

	// VALIDAZIONE HEADER GIF E LUNGHEZZA
	isGIF, err := regexp.MatchString(`^GIF`, string(u.Photo))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !isGIF || len(u.Photo) < 4 || len(u.Photo) > 524288 {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	// SALVO IN DATABASE
	err = rt.db.SetMyUserphoto(u.Photo, token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// OK
	w.WriteHeader(http.StatusNoContent)
}
