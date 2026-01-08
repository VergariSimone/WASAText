package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	token := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")

	convid, err := strconv.Atoi(ps.ByName("convid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	mesid, err := strconv.Atoi(ps.ByName("mesid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	n, err := rt.db.DeleteMessage(token, mesid, convid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if n < 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if n == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
