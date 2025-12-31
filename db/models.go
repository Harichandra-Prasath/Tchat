package db

type UserModel struct {
	Username string `db:"username"`
	Password string `db:"password"`
}
