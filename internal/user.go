package internal

import (
	"database/sql"
	"errors"
	"observe/schema"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *sql.DB, user schema.User) (schema.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return schema.User{}, errors.New("Error hashing password: " + err.Error())
	}

	query := `
    INSERT INTO users (username, password, created_at, updated_at)
    VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    RETURNING id, username, password, registered_at;
  `
	_, err = db.Exec(query, user.Username, string(hashedPassword))
	if err != nil {
		return schema.User{}, errors.New("Error querying database: " + err.Error())
	}
	return user, nil
}

func GetAllUsers(db *sql.DB) ([]schema.User, error) {
	query := `
        SELECT * FROM users;
  `
	rows, err := db.Query(query)
	if err != nil {
		return nil, errors.New("Error querying database: " + err.Error())
	}
	defer rows.Close()

	var users []schema.User
	for rows.Next() {
		var user schema.User
		// if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.RegisteredAt); err != nil {
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, errors.New("Error scanning rows: " + err.Error())
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.New("Error reading rows: " + err.Error())
	}
	return users, nil
}

func GetUserByID(db *sql.DB, userID int) (schema.User, error) {
	query := `
    SELECT id, username, password, created_at, updated_at FROM users WHERE id = $1;
  `
	var user schema.User
	err := db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return schema.User{}, errors.New("user not found")
		}
		return schema.User{}, errors.New("Error querying database: " + err.Error())
	}
	return user, nil
}

func GetUserByUsername(db *sql.DB, username string) (schema.User, error) {
	query := `
    SELECT id, username, password, created_at, updated_at FROM users WHERE username = $1;
  `
	var user schema.User
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return schema.User{}, errors.New("user not found")
		}
		return schema.User{}, errors.New("Error querying database: " + err.Error())
	}
	return user, nil
}

func UpdateUser(db *sql.DB, user schema.User) (schema.User, error) {
	query := `
    UPDATE users
    SET username = $1, password = $2, updated_at = CURRENT_TIMESTAMP
    WHERE id = $3
    RETURNING id, username, password, created_at, updated_at;
  `
	err := db.QueryRow(query, user.Username, user.Password, user.ID).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return schema.User{}, errors.New("Error querying database: " + err.Error())
	}
	return user, nil
}

func DeleteUser(db *sql.DB, userID int) error {
	query := `
    DELETE FROM users WHERE id = $1;
  `
	result, err := db.Exec(query, userID)
	if err != nil {
		return errors.New("Error executing database query: " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("Error getting rows affected: " + err.Error())
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func VerifyUser(user schema.User, db *sql.DB) error {
	userFromDB, err := GetUserByUsername(db, user.Username)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(user.Password))
	if err != nil {
		return errors.New("invalid password")
	}
	return nil
}
