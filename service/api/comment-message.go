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

	// 1. Parse ID Messaggio
	mesID, err := strconv.Atoi(ps.ByName("mesid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 2. Parse Token
	token := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")

	// 3. Parse Body (JSON)
	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 4. Chiama il DB
	id, err := rt.db.CommentMessage(mesID, token, req.Emoji)
	if err != nil {
		// Se l'errore è SQL (es. utente non trovato), diamo 400
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 5. Rispondi OK
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    id,
		"emoji": req.Emoji,
	})
}
