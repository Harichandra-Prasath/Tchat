package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/Harichandra-Prasath/Tchat/configs"
	"github.com/jackc/pgx/v5"
)

func IsUserExists(username string) (bool, error) {

	var exists bool

	existsQuery := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE username=$1)`, configs.TbCfg.UsersTable)
	err := DBPool.QueryRow(context.Background(), existsQuery, username).Scan(&exists)
	if err != nil {
		return exists, err
	}

	return exists, nil

}

func GetUser(username string) (*UserModel, error) {

	var u UserModel
	query := fmt.Sprintf(`SELECT id,username,password FROM %s WHERE username = $1`, configs.TbCfg.UsersTable)
	err := DBPool.QueryRow(context.Background(), query, username).Scan(&u.Id, &u.Username, &u.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, UserDoesNotExistsError
		}
		return nil, err
	}

	return &u, nil

}
