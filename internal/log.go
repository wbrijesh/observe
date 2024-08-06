package internal

import (
	"database/sql"
	"errors"
	"observe/schema"
	"time"
)

func InsertLog(db *sql.DB, log schema.Log) (schema.Log, error) {
	query := `
    INSERT INTO logs (project_id, message, level, timestamp)
    VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
    RETURNING id, timestamp;
  `
	err := db.QueryRow(query, log.ProjectID, log.Message, log.Level).Scan(&log.ID, &log.Timestamp)
	if err != nil {
		return schema.Log{}, errors.New("Error inserting log: " + err.Error())
	}
	return log, nil
}

func BatchInsertLogs(db *sql.DB, logs []schema.Log) ([]schema.Log, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.New("Error starting transaction: " + err.Error())
	}

	query := `
    INSERT INTO logs (project_id, message, level, timestamp)
    VALUES ($1, $2, $3, $4)
    RETURNING id;
  `

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("Error preparing statement: " + err.Error())
	}
	defer stmt.Close()

	for i := range logs {
		err := stmt.QueryRow(logs[i].ProjectID, logs[i].Message, logs[i].Level, logs[i].Timestamp).Scan(&logs[i].ID)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("Error inserting log: " + err.Error())
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.New("Error committing transaction: " + err.Error())
	}
	return logs, nil
}

func DeleteLogsByTimeRange(db *sql.DB, startTime, endTime time.Time) error {
	query := `
    DELETE FROM logs
    WHERE timestamp BETWEEN $1 AND $2;
  `
	_, err := db.Exec(query, startTime, endTime)
	if err != nil {
		return errors.New("Error deleting logs by time range: " + err.Error())
	}
	return nil
}

func DeleteLogsByProject(db *sql.DB, projectID int) error {
	query := `
    DELETE FROM logs
    WHERE project_id = $1;
  `
	_, err := db.Exec(query, projectID)
	if err != nil {
		return errors.New("Error deleting logs by project: " + err.Error())
	}
	return nil
}
