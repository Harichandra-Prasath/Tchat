package db

import (
	"context"
	"fmt"

	"github.com/Harichandra-Prasath/Tchat/configs"
	"github.com/google/uuid"
)

func CreateUser(user *UserModel) (uuid.UUID, error) {

	var id uuid.UUID
	var exists bool

	existsQuery := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE username=$1)`, configs.TbCfg.UsersTable)
	err := DBPool.QueryRow(context.Background(), existsQuery, user.Username).Scan(&exists)
	if err != nil {
		return uuid.Nil, fmt.Errorf("checking user: %s", err.Error())
	}
	if exists {
		return uuid.Nil, UserExistsError
	}

	insertQuery := fmt.Sprintf(`INSERT INTO %s (username, password) VALUES ($1, $2) RETURNING id`, configs.TbCfg.UsersTable)

	err = DBPool.QueryRow(context.Background(), insertQuery, user.Username, user.Password).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("inserting user: %s", err.Error())
	}

	return id, nil
}
