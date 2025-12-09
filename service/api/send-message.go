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

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("content-type", "application/json")

	// ESTRAZIONE TOKEN
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	conversationParam := r.URL.Query().Get("conversation")
	replyParam := r.URL.Query().Get("reply")

	// VALIDAZIONE PARAMENTRO CONVERSATION
	if conversationParam == "" {
		w.WriteHeader(http.StatusNotAcceptable) // 406 se non è presente
		return
	}

	conID, msgID, err := rt.resolveConversationAndReply(token, conversationParam, replyParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// PARSING DEL PAYLOAD
	text, image, err := rt.parseMessagePayload(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// VALIDAZIONE CONTENT + CATEGORIA
	category, err := validateMessageContent(text, image)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable) // 406 se testo o immagine non rispettano i limiti
		return
	}

	// INVIO MESS AL DB
	if err := rt.db.SendMessage(conID, token, text, category, image, msgID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// RECUPERO CONV AGGIORNATA
	conv, err := rt.db.GetConversation(conID, 0, 50, token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(conv)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated) // 201
	if _, err := w.Write(data); err != nil {
		rt.baseLogger.WithError(err).Warning("Errore durante la scrittura della risposta")
	}
}

// FUNZIONI DI SUPPORT

// resolveConversationAndReply
// RECUPERO/ CREAZIONE CONVERSAZIONE , VALIDAZIONE TOKEN E GESTIONE REPLY
func (rt *_router) resolveConversationAndReply(token, conversationParam, replyParam string) (int, int, error) {
	isNumeric := regexp.MustCompile(`^\d+$`)

	var conID int
	var err error

	//SE CONV È NUMERICA, ALLORA ID
	if isNumeric.MatchString(conversationParam) {
		conID, err = strconv.Atoi(conversationParam)
		if err != nil {
			return 0, -1, err
		}
		// AUTORIZZAZIONE
		if err := rt.db.IsAuthorized(token, conID); err != nil {
			return 0, -1, err
		}
	} else {
		// ALTRIMENTI È USERNAME -> IMPOSTA CONVERSAZIONE
		conID, err = rt.db.SetConversation(token, "conversation", conversationParam, "", nil)
		if err != nil {
			return 0, -1, err
		}
	}

	// GESTIONE REPLU
	msgID := -1
	if isNumeric.MatchString(replyParam) {
		msgID, err = strconv.Atoi(replyParam)
		if err != nil {
			return 0, -1, err
		}
		if err := rt.db.IsMsgAuthorized(token, msgID); err != nil {
			return 0, -1, err
		}
	}

	return conID, msgID, nil
}

// parseMessagePayload
// LEGGE DATI INVIATI DAL CLIENT
func (rt *_router) parseMessagePayload(r *http.Request) (string, []byte, error) {
	imgFile, _, err := r.FormFile("image")
	var image []byte

	if err == nil {
		defer imgFile.Close()
		image, err = io.ReadAll(imgFile)
		if err != nil {
			return "", nil, err
		}
	}

	text := r.FormValue("text")

	return text, image, nil
}

// validateMessageContent
// VALIDAZIONE TESTO E MESSAGGI , ASSEGNA CATEGORIA AL TIPO DI MESSAGGIOS
func validateMessageContent(text string, image []byte) (string, error) {
	hasText := len(text) > 0
	hasImage := len(image) > 0

	if hasText && (len(text) < 1 || len(text) > 8192) {
		return "", http.ErrBodyNotAllowed
	}

	if hasImage && (len(image) < 64 || len(image) > 524288) {
		return "", http.ErrBodyNotAllowed
	}

	switch {
	case hasText && hasImage:
		return "mixed", nil
	case hasText:
		return "text", nil
	default:
		return "image", nil
	}
}
