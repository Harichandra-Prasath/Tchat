package db

import (
	"context"
	"fmt"

	"github.com/Harichandra-Prasath/Tchat/configs"
	"github.com/google/uuid"
)

func CreateUser(user *UserModel) (uuid.UUID, error) {

	var id uuid.UUID

	exists, err := IsUserExists(user.Username)
	if err != nil {
		return uuid.Nil, fmt.Errorf("checking user existense: %s", err.Error())
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

func DeleteSession(userId uuid.UUID) error {
	deleteQuery := fmt.Sprintf(`DELETE FROM %s WHERE user_id=$1`, configs.TbCfg.SessionsTable)
	_, err := DBPool.Exec(context.Background(), deleteQuery, userId)
	if err != nil {
		return fmt.Errorf("deleting user session: %s", err.Error())
	}
	return nil
}

func CreateSession(session *SessionModel) error {
	writeQuery := fmt.Sprintf(`INSERT INTO %s (token,user_id,created_at,expires_at) VALUES ($1, $2, $3, $4)`, configs.TbCfg.SessionsTable)
	_, err := DBPool.Exec(context.Background(), writeQuery, session.Token, session.UserId, session.CreatedAt, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("inserting session: %s", err.Error())
	}

	return nil

}
