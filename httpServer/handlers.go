package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Harichandra-Prasath/Tchat/db"
	"github.com/Harichandra-Prasath/Tchat/logging"
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
		id, err := db.CreateUser(&userModel)
		if err != nil {

			if errors.Is(err, db.UserExistsError) {
				logging.Logger.Info("Username Exists Already")
				http.Error(w, err.Error(), 400)
				return
			}

			logging.Logger.Error("Error in User Creation", "err", err.Error())
			http.Error(w, "User Creation Failed", 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)

		err = json.NewEncoder(w).Encode(registerUserResponseSchema{Id: id, Message: "User Created Successfully"})
		if err != nil {
			logging.Logger.Error("Error in User Creation", "err", err.Error())
			http.Error(w, "User Createion Failed", 500)
			return
		}
	})
}
