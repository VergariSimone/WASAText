package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"regexp"
)

// STRUTTURE
type Name struct {
	Name string `json:"name"`
}

type Session struct {
	Identifier string `json:"identifier"`
}

// ESPRESSIONE REGEX
var validName = regexp.MustCompile(`^[A-Za-z].*$`)

// HANDLER
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// DECODIFICA JSON
	var n Name
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, `{"error": "json non valido"}`, http.StatusBadRequest)
		return
	}

	// VALIDAZIONE NOME
	if !validName.MatchString(n.Name) {
		http.Error(w, `{"error": "nome non valido: deve contenere almeno una lettera"}`, http.StatusBadRequest)
		return
	}

	// GENERAZIONE TOKEN E SESSIONE
	identifier, err := rt.db.SetToken(n.Name)
	if err != nil {
		http.Error(w, `{"error": "errore interno del server"}`, http.StatusInternalServerError)
		rt.baseLogger.WithError(err).Error("Errore SetToken")
		return
	}

	// JSON RISPOSTA
	session := Session{Identifier: identifier}

	w.WriteHeader(http.StatusCreated) // 201
	if err := json.NewEncoder(w).Encode(session); err != nil {
		rt.baseLogger.WithError(err).Warning("Errore encoding JSON verso http.ResponseWriter")
	}
}
