package gateway

import (
	"backend/src/domain"
	"backend/src/infra/db"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// user作成
func CreateUser(user domain.UserInput) (userId int, err error) {
	const MAX_USER_NAME int = 50
	if len(user.UserName) > MAX_USER_NAME {
		log.Printf("error: User name too long.")
		return 0, fmt.Errorf(domain.BadRequest)
	}

	_, _, err = GetUser(user.UserName)
	if err == nil {
		log.Printf("error: User name %s already exist.", user.UserName)
		return 0, fmt.Errorf(domain.Conflict)
	}

	db := db.OpenDB()
	defer db.Close()

	queryCreateUser := "insert into users (user_name, password) values (?, ?);"

	result, err := db.Exec(queryCreateUser, user.UserName, user.Password)
	if err != nil {
		log.Println("error: ", err)
		return 0, fmt.Errorf(domain.InternalServerError)
	}

	// LastInsertId()で今挿入したuser_idを取得
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Println("error: ", err)
		return 0, fmt.Errorf(domain.InternalServerError)
	}

	return int(lastInsertId), nil
}

// user取得
func GetUser(userName string) (userId int, hashedPassword string, err error) {
	db := db.OpenDB()
	defer db.Close()

	queryGetUser := "select user_id, password from users where user_name = ?;"
	userRows := db.QueryRow(queryGetUser, userName)

	if err := userRows.Scan(&userId, &hashedPassword); err != nil {
		log.Println("error: ", err)

		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", fmt.Errorf(domain.NotFound)
		}

		return 0, "", fmt.Errorf(domain.InternalServerError)
	}

	return userId, hashedPassword, nil
}

func DeleteUserByUserId(userId int) error {
	db := db.OpenDB()
	defer db.Close()

	queryDeleteUser := "delete from users where user_id = ?;"
	res, err := db.Exec(queryDeleteUser, userId)
	if err != nil {
		log.Println("error: ", err)
		return fmt.Errorf(domain.InternalServerError)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("error: ", err)
		return fmt.Errorf(domain.InternalServerError)
	}
	if rowsAffected == 0 {
		log.Printf("no user found with user_id %d", userId)
		return fmt.Errorf(domain.NotFound)
	}

	return nil
}
