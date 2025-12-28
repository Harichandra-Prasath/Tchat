package httpserver

type schema any

type messageMetadata struct{}

type message struct {
	Sender   string          `json:"sender" validate:"required"`
	Metadata messageMetadata `json:"metadata"`
	Message  string          `json:"message" validate:"required"`
	Reciever string          `json:"reciever" validate:"required"`
}
