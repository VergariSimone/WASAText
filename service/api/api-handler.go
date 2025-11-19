package api

import (
	"net/http"
)

// // Handler restituisce un'istanza di httprouter.Router che gestisce le API registrate qui
func (rt *_router) Handler() http.Handler {

	// Registra routes

	// lOGIN
	rt.router.POST("/session",rt.doLogin)


	// INFO PROFILO UTENTE
	rt.router.GET("/info", rt.getUserInfo)
	rt.router.PUT("/info/username", rt.setMyUserName)
	rt.router.PUT("/info/userphoto", rt.setMyPhoto)

	// CONVERSAZIONI GENERALI
	rt.router.GET("/conversations", rt.getMyConversations)
	rt.router.POST("/conversations", rt.sendMessage)
	rt.router.PUT("/conversations/group", rt.addGroup)

	// CONVERSAZIONI SPECIFICHE
	rt.router.GET("/conversations/byid/:convid", rt.getConversations)
	rt.router.PUT("/conversations/byid/:convid", rt.leaveGroup)
	rt.router.PUT("/conversations/byid/:convid/name", rt.setGroupName)
	rt.router.PUT("/conversations/byid/:convid/photo", rt.setGroupPhoto)

	// INTERRAZIONE MESSAGGI SPECIFICI
	rt.router.POST("/conversations/byid/:convid/message/:mesid", rt.forwardMessage)
	rt.router.DELETE("/conversations/byid/:convid/message/:mesid", rt.deleteMessage)

	// COMMENTI 
	rt.router.POST("/conversations/byid/:convid/message/:mesid/comment", rt.commentMessage)
	rt.router.DELETE("/conversations/byid/:convid/message/:mesid/comment/:commid", rt.uncommentMessage)

	// Routes speciali
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
