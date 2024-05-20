//Data access layer responsible for interacting with the database.

package authDao

import (
	"context"
	"database/sql"
	"log"

	"github.com/naman2607/netflixClone/database"
	user "github.com/naman2607/netflixClone/models"
)

func CheckIfUserExists(email string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	var exists bool

	err := database.ExecuteTransactional(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, query, email)
		if err := row.Scan(&exists); err != nil {
			if err == sql.ErrNoRows {
				exists = false
				return nil
			}
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	return exists, nil
}

func InsertUserIntoDatabase(user user.UserBasicDetails) error {
	err := database.ExecuteTransactional(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, "INSERT INTO users (username,email,password) VALUES ($1 , $2, $3)", user.Username, user.Email, user.Password)
		if err != nil {
			log.Println("Error in executing insert user query ", err)
			return err
		}
		log.Println("User details saved to db")
		return nil
	})
	return err
}

func GetUser(email string) (user.UserBasicDetails, error) {
	var user user.UserBasicDetails
	err := database.ExecuteTransactional(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		query := "SELECT username, email, password FROM users WHERE email = $1"
		rows := tx.QueryRowContext(ctx, query, email)
		err := rows.Scan(&user.Username, &user.Email, &user.Password)
		if err != nil {
			log.Println("eror in scanning get user row", err)
		}
		return nil
	})
	return user, err
}
