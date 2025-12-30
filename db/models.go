package db

import "github.com/google/uuid"

type UserModel struct {
	Id       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Password string    `db:"password"`
}
