package schema

import "time"

type User struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
}

type Project struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Enviroment string    `json:"enviroment"`
	Name       string    `json:"name"`
}

type Log struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Level     string    `json:"level"`
}
