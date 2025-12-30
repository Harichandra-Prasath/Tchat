package httpserver

import (
	"net/http"

	"github.com/Harichandra-Prasath/Tchat/db"
	"github.com/Harichandra-Prasath/Tchat/utils"
)

func sendMessageHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func sseHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func registerHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var key ctxKey[registerUserSchema]
		reqData := r.Context().Value(key).(registerUserSchema)

		userModel := db.UserModel{Username: reqData.Username, Password: utils.HashPassword(reqData.Password)}

	})
}
