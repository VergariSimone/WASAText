package api

/*parameters:
	- $ref: "#/components/parameters/ConversationID"
	- $ref: "#/components/parameters/MessageID"
  requestBody:
	description: reaction emoji
	required: true
	content:
	  image/gif:
		schema:
		  description: emoji
		  type: string
		  format: binary
		  minLength: 64
		  maxLength: 524288
  responses:
	"200":
	  description: message commented
	  content:
		application/json:
		  schema:
			$ref: "#/components/schemas/ComID"
	"401":
	  $ref: "#/components/responses/UnauthorizedError"
	"404":
	  description: message id or conversation id not found
*/

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	responseData, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// fmt.Println(len(responseData), len(string(responseData)))
	if (len(responseData) > 4) || (len(responseData) < 1) {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	emoji := string(responseData)

	err = rt.db.IsAuthorized(token, convid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}


	_, err = rt.db.CommentMessage(mesid, token, emoji)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	
}