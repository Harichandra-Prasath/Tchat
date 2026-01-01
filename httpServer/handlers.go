package httpserver

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

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

func loginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var key ctxKey[loginSchema]
		reqData := r.Context().Value(key).(loginSchema)

		user, err := db.GetUser(reqData.Username)
		if err != nil {
			if errors.Is(err, db.UserDoesNotExistsError) {
				logging.Logger.Info("User Doesnt Exist")
				http.Error(w, "Invalid Credentials", 400)
				return
			}
			logging.Logger.Error("Error in Getting the User", "err", err.Error())
			http.Error(w, "Login Failed", 500)
		}

		if !utils.VerifyPassword(user.Password, reqData.Password) {
			http.Error(w, "Invalid Credentials", 400)
			return
		}

		raw := make([]byte, 32)
		rand.Read(raw)

		token := base64.RawURLEncoding.EncodeToString(raw)

		hashedToken := utils.HashToken(token)

		var s db.SessionModel
		s.UserId = user.Id
		s.Token = hashedToken
		s.CreatedAt = time.Now()
		s.ExpiresAt = s.CreatedAt.Add(30 * 24 * time.Hour)

		err = db.CreateSession(&s)
		if err != nil {
			logging.Logger.Error("Error in Login", "err", err.Error())
			http.Error(w, "Login Failed", 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		err = json.NewEncoder(w).Encode(loginResponseSchema{Token: token, Message: "Login Succesfull"})
		if err != nil {
			logging.Logger.Error("Error in Login", "err", err.Error())
			http.Error(w, "Login Failed", 500)
			return
		}

	})

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
