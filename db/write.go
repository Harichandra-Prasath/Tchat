package db

import "fmt"

func CreateUser(user *UserModel) error {

	query := fmt.Sprintf(`INSERT INTO %s (id, username, password) VALUES ($1, $2)`)
}
