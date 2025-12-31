package httpserver

import "github.com/google/uuid"

type messageMetadata struct{}

type sendMessageSchema struct {
	Sender   string          `json:"sender" validate:"required"`
	Metadata messageMetadata `json:"metadata"`
	Message  string          `json:"message" validate:"required"`
	Reciever string          `json:"reciever" validate:"required"`
}

type registerUserSchema struct {
	Username string `json:"username" validate:"required,min=3,max=12"`
	Password string `json:"password" validate:"required,min=8,max=16"`
}

type loginSchema struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type registerUserResponseSchema struct {
	Id      uuid.UUID `json:"id"`
	Message string    `json:"message"`
}

type loginResponseSchema struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}
