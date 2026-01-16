package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type CommentRequest struct {
	Emoji string `json:"text"`
}

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// parse id messaggio
	mesID, err := strconv.Atoi(ps.ByName("mesid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// parse token
	token := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")

	// parse body (JSON)
	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// chiama il DB
	id, err := rt.db.CommentMessage(mesID, token, req.Emoji)
	if err != nil {
		// se l'errore è SQL (tipo: utente non trovato), diamo 400
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// rispondi OK
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    id,
		"emoji": req.Emoji,
	})
	if err != nil {
		_ = err

	}
}
