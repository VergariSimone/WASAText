package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

// // Handler restituisce un'istanza di httprouter.Router che gestisce le API registrate qui
func (rt *_router) Handler() http.Handler {

	// Registra routes

	// lOGIN
	rt.router.POST("/session", rt.doLogin)

	// INFO PROFILO UTENTE
	rt.router.GET("/info", rt.authControl(rt.getUserInfo))
	rt.router.GET("/users", rt.authControl(rt.getUsers))
	rt.router.PUT("/info/username", rt.authControl(rt.setMyUserName))
	rt.router.PUT("/info/userphoto", rt.authControl(rt.setMyPhoto))

	// CONVERSAZIONI GENERALI
	rt.router.GET("/conversations", rt.authControl(rt.getMyConversations))
	rt.router.POST("/conversations", rt.authControl(rt.sendMessage))
	rt.router.PUT("/conversations/group", rt.authControl(rt.addToGroup))

	// CONVERSAZIONI SPECIFICHE
	rt.router.GET("/conversations/byid/:convid", rt.authControl(rt.getConversation))
	rt.router.PUT("/conversations/byid/:convid", rt.authControl(rt.leaveGroup))
	rt.router.PUT("/conversations/byid/:convid/name", rt.authControl(rt.setGroupName))
	rt.router.PUT("/conversations/byid/:convid/photo", rt.authControl(rt.setGroupPhoto))
	rt.router.GET("/conversations/byid/:convid/info", rt.authControl(rt.getGroupInfo))

	// INTERRAZIONE MESSAGGI SPECIFICI
	rt.router.POST("/conversations/byid/:convid/message/:mesid", rt.authControl(rt.forwardMessage))
	rt.router.DELETE("/conversations/byid/:convid/message/:mesid", rt.authControl(rt.deleteMessage))

	// COMMENTI
	rt.router.POST("/conversations/byid/:convid/message/:mesid/comment", rt.authControl(rt.commentMessage))
	rt.router.DELETE("/conversations/byid/:convid/message/:mesid/comment/:commid", rt.authControl(rt.uncommentMessage))

	// Routes speciali
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}

func (rt *_router) authControl(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// SE TOKEN VALIDO RECUOERA NOME UTENRE
		if name, err := rt.db.GetName(token); err != nil || name == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if _, err := rt.db.GetToken(token); err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next(w, r, ps)
	}
}
