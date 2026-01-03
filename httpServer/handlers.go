package httpserver

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Harichandra-Prasath/Tchat/broker"
	"github.com/Harichandra-Prasath/Tchat/db"
	"github.com/Harichandra-Prasath/Tchat/logging"
	"github.com/Harichandra-Prasath/Tchat/utils"
)

func sendMessageHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userKey ctxUser
		sender := r.Context().Value(userKey).(*db.UserModel)

		var dataKey ctxKey[sendMessageSchema]
		reqData := r.Context().Value(dataKey).(sendMessageSchema)

		// check for the reciever
		reciever, err := db.GetUser(reqData.Reciever)
		if err != nil {
			if errors.Is(err, db.UserDoesNotExistsError) {
				logging.Logger.Info("User Doesnt Exist")
				http.Error(w, "Invalid Reciever", 400)
				return
			}
			logging.Logger.Error("Error in Getting the Reciever", "err", err.Error())
			http.Error(w, "Send message Failed", 500)
			return
		}

		m := message{Sender: sender.Username, Message: reqData.Message}

		if ch, ok := ChnMapper.Load(reciever.Id); !ok {
			logging.Logger.Info("Reciever not Online. Publishing to transient storage")
			data, _ := json.Marshal(m)
			err := broker.PublishEvents(&broker.Event{Type: "message", DSN: reciever.Username, Data: data})
			if err != nil {
				logging.Logger.Error("Publishing Message", "err", err.Error())
				http.Error(w, "Send Message Failed", 500)
				return
			}
			logging.Logger.Info("Message Published")
		} else {
			logging.Logger.Info("Reciever Online. Forwarding the Message")
			ch.(chan *message) <- &m
		}

		w.WriteHeader(201)
		w.Write([]byte("Message Produced"))
	})
}

func sseHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var key ctxUser
		user := r.Context().Value(key).(*db.UserModel)

		chn := getChannel(user.Id)

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("connection", "keep-alive")
		w.Header().Set("cache-control", "no-cache")

		flusher, _ := w.(http.Flusher)

		chn <- &message{Sender: "system", Message: "You are Connected"}

		eventChan := make(chan *broker.Event)
		go func() {
			err := broker.ConsumeEvents(eventChan, user.Username)
			if err != nil {
				logging.Logger.Error("Consuming Events", "err", err.Error())
				return
			}

		}()

		for {
			select {
			case <-r.Context().Done():
				// Cleanup
				logging.Logger.Warn("Client Disconnected.", "usr", user.Username)
				ChnMapper.Delete(user.Id)
				return

			case m := <-chn:
				data, _ := json.Marshal(m)
				fmt.Fprintf(w, "event: message\n")
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()

			case e := <-eventChan:
				fmt.Fprintf(w, "event:%s\n", e.Type)
				fmt.Fprintf(w, "data: %s\n\n", e.Data)
				flusher.Flush()
			}
		}

	})
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
			return
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

		err = db.DeleteSession(user.Id)
		if err != nil {
			logging.Logger.Error("Error in Login", "err", err.Error())
			http.Error(w, "Login Failed", 500)
			return
		}

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

func logoutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var key ctxUser
		user := r.Context().Value(key).(*db.UserModel)
		err := db.DeleteSession(user.Id)
		if err != nil {
			logging.Logger.Error("Error in Logout", "err", err.Error())
			http.Error(w, "Logout Failed", 500)
			return
		}
		w.WriteHeader(200)
		w.Write([]byte("Logout Successfull"))
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
