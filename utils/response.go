package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"observe/schema"
)

func SendResponse(w http.ResponseWriter, r *http.Request, response schema.Response) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleError(w http.ResponseWriter, r *http.Request, statusCode int, message string, err error) {
	if err != nil {
		w.WriteHeader(statusCode)
		response := schema.Response{
			Status:  "ERROR",
			Message: message + err.Error(),
		}
		SendResponse(w, r, response)
	}
}

func HandleMethodNotAllowed(w http.ResponseWriter, r *http.Request, method string) {
	if r.Method != method {
		HandleError(w, r, http.StatusMethodNotAllowed, "", errors.New("method"+method+"not allowed"))
	}
}
