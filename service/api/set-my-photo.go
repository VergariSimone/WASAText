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

	// CONTROLLO CONTENT-TYPE (GIF, PNG o JPEG)
	ct := r.Header.Get("Content-Type")
	if ct != "image/gif" && ct != "image/png" && ct != "image/jpeg" {
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

	// VALIDAZIONE FORMATO 
	isGIF, _ := regexp.MatchString(`^GIF`, string(u.Photo))
	isPNG, _ := regexp.MatchString(`^\x89PNG`, string(u.Photo))
	isJPG, _ := regexp.MatchString(`^\xff\xd8`, string(u.Photo))

	// se non è nessuno dei tre formati, o se dimensioni errate -> 406
	if (!isGIF && !isPNG && !isJPG) || len(u.Photo) < 4 || len(u.Photo) > 1048576 {
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