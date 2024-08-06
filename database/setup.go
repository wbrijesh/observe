package database

import (
	"database/sql"
)

func CreateUsersTable(db *sql.DB) {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS users (
      id VARCHAR(255) PRIMARY KEY,
      username VARCHAR(255) NOT NULL,
      password VARCHAR(255) NOT NULL,
      created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );
  `)
	if err != nil {
		panic(err)
	}
}

func CreateProjectsTable(db *sql.DB) {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS projects (
      id VARCHAR(255) PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      environment VARCHAR(255) NOT NULL,
      user_id VARCHAR(255) NOT NULL,
      created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (user_id) REFERENCES users (id)
    );
  `)
	if err != nil {
		panic(err)
	}
}

func CreateLogsTable(db *sql.DB) {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS logs (
      id VARCHAR(255) PRIMARY KEY,
      project_id VARCHAR(255) NOT NULL,
      message TEXT NOT NULL,
      level VARCHAR(255) NOT NULL,
      timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,  -- Log timestamp
      FOREIGN KEY (project_id) REFERENCES projects (id)
    );
  `)
	if err != nil {
		panic(err)
	}
}

func CreateIndexes(db *sql.DB) {
	_, err := db.Exec(`
  CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
  CREATE INDEX IF NOT EXISTS idx_logs_project_id ON logs(project_id);
  CREATE INDEX IF NOT EXISTS idx_logs_level ON logs(level);
  `)
	if err != nil {
		panic(err)
	}
}
