package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"observe/internal"
	"observe/schema"
	"observe/utils"
	"observe/validation"
)

func UserRegistrationHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	utils.HandleMethodNotAllowed(w, r, http.MethodPost)

	var user schema.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.HandleError(w, r, http.StatusBadRequest, "Invalid request body: ", err)
		return
	}
	user.ID = utils.GenerateUUID()

	err = validation.ValidateUserForRegistration(user)
	if err != nil {
		utils.HandleError(w, r, http.StatusBadRequest, "Invalid user data: ", err)
		return
	}

	_, err = internal.GetUserByUsername(db, user.Username)
	if err == nil {
		utils.HandleError(w, r, http.StatusConflict, "", errors.New("user already exists"))
		return
	}

	_, err = internal.CreateUser(db, user)
	if err != nil {
		utils.HandleError(w, r, http.StatusInternalServerError, "Failed to create user: ", err)
		return
	}

	token, err := internal.GenerateToken(user.Username)
	if err != nil {
		utils.HandleError(w, r, http.StatusInternalServerError, "Failed to generate token: ", err)
		return
	}

	response := schema.Response{
		Status:  "SUCCESS",
		Message: "User created successfully",
		Data: map[string]string{
			"token": token,
		},
	}

	utils.SendResponse(w, r, response)
}

func UserAssertionHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	utils.HandleMethodNotAllowed(w, r, http.MethodPost)
	var user schema.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.HandleError(w, r, http.StatusBadRequest, "Invalid request body: ", err)
		return
	}

	err = validation.ValidateUserForLogin(user)
	if err != nil {
		utils.HandleError(w, r, http.StatusBadRequest, "Invalid user data: ", err)
		return
	}

	err = internal.VerifyUser(user, db)
	if err != nil {
		utils.HandleError(w, r, http.StatusUnauthorized, "Invalid email or password", err)
		return
	}

	token, err := internal.GenerateToken(user.Username)
	if err != nil {
		utils.HandleError(w, r, http.StatusInternalServerError, "Failed to generate token: ", err)
		return
	}
	response := schema.Response{
		Status:  "SUCCESS",
		Message: "User logged in successfully",
		Data:    map[string]string{"token": token},
	}
	utils.SendResponse(w, r, response)
}

// write CURL requests to test the handlers at localhost:8080 and /register and /login endpoints
// curl -X POST -H "Content-Type: application/json" -d '{"username": "test", "password": "test"}' http://localhost:8080/register
// curl -X POST -H "Content-Type: application/json" -d '{"username": "test", "password": "test"}' http://localhost:8080/login
