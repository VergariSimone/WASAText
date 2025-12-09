package api

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type User struct {
	Name  string `json:"name"`
	Photo string `json:"photo"` // restituita in formato binhex
}

func (rt *_router) getUserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// ESTRAZIONE TOKEN 
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// RECUPERO INFO UTENTE
	name, photoBytes, err := rt.db.GetUserInfo(token)
	if err != nil {
		name = ""
		photoBytes = nil
	}

	// CONVERSIONE IN  BINHEX (hexadecimal ASCII)
	photoHex := ""
	if photoBytes != nil {
		photoHex = hex.EncodeToString(photoBytes)
	}

	userResp := User{
		Name:  name,
		Photo: photoHex,
	}

	resp, err := json.Marshal(userResp)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		rt.baseLogger.WithError(err).Warning("Errore nella scrittura della risposta HTTP")
	}
}
