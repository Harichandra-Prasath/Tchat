package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/Harichandra-Prasath/Tchat/configs"
	"github.com/Harichandra-Prasath/Tchat/utils"
	"github.com/google/uuid"
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

func GetUserbyID(id uuid.UUID) (*UserModel, error) {

	var u UserModel
	query := fmt.Sprintf(`SELECT id,username FROM %s WHERE id = $1`, configs.TbCfg.UsersTable)
	err := DBPool.QueryRow(context.Background(), query, id).Scan(&u.Id, &u.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, UserDoesNotExistsError
		}
		return nil, err
	}

	return &u, nil

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

func GetSession(token string) (uuid.UUID, error) {

	tokenHash := utils.HashToken(token)
	var userId uuid.UUID
	query := fmt.Sprintf(`SELECT user_id FROM %s WHERE token = $1 AND expires_at>now()`, configs.TbCfg.SessionsTable)

	err := DBPool.QueryRow(context.Background(), query, tokenHash).Scan(&userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, SessionDoesNotExistsError
		}

		return uuid.Nil, err
	}

	return userId, nil
}
