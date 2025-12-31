package db

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	Id       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Password string    `db:"password"`
}

type SessionModel struct {
	Token     string    `db:"token"`
	UserId    uuid.UUID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}
