package internal

import (
	"database/sql"
	"errors"
	"observe/schema"
)

func CreateProject(db *sql.DB, project schema.Project) (schema.Project, error) {
	query := `
    INSERT INTO projects (name, enviroment, user_id, created_at, updated_at)
    VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    RETURNING id, created_at, updated_at;
  `
	err := db.QueryRow(query, project.Name, project.Enviroment, project.UserID).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		return schema.Project{}, errors.New("Error creating project: " + err.Error())
	}
	return project, nil
}

func GetAllProjects(db *sql.DB) ([]schema.Project, error) {
	query := `
    SELECT id, name, enviroment, user_id, created_at, updated_at FROM projects;
  `
	rows, err := db.Query(query)
	if err != nil {
		return nil, errors.New("Error querying all projects: " + err.Error())
	}
	defer rows.Close()

	var projects []schema.Project
	for rows.Next() {
		var project schema.Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Enviroment, &project.UserID, &project.CreatedAt, &project.UpdatedAt); err != nil {
			return nil, errors.New("Error scanning project: " + err.Error())
		}
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.New("Error iterating over projects: " + err.Error())
	}
	return projects, nil
}

func GetProjectsByUserID(db *sql.DB, userID int) ([]schema.Project, error) {
	query := `
    SELECT id, name, enviroment, user_id, created_at, updated_at FROM projects WHERE user_id = $1;
  `
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, errors.New("Error querying projects by user ID: " + err.Error())
	}
	defer rows.Close()

	var projects []schema.Project
	for rows.Next() {
		var project schema.Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Enviroment, &project.UserID, &project.CreatedAt, &project.UpdatedAt); err != nil {
			return nil, errors.New("Error scanning project: " + err.Error())
		}
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.New("Error iterating over projects: " + err.Error())
	}
	return projects, nil
}

func GetProjectByID(db *sql.DB, projectID int) (schema.Project, error) {
	query := `
    SELECT id, name, enviroment, user_id, created_at, updated_at FROM projects WHERE id = $1;
  `
	var project schema.Project
	err := db.QueryRow(query, projectID).Scan(&project.ID, &project.Name, &project.Enviroment, &project.UserID, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return schema.Project{}, errors.New("project not found")
		}
		return schema.Project{}, errors.New("Error querying project by ID: " + err.Error())
	}
	return project, nil
}

func UpdateProject(db *sql.DB, project schema.Project) (schema.Project, error) {
	query := `
    UPDATE projects
    SET name = $1, enviroment = $2, updated_at = CURRENT_TIMESTAMP
    WHERE id = $3
    RETURNING created_at, updated_at;
  `
	err := db.QueryRow(query, project.Name, project.Enviroment, project.ID).Scan(&project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		return schema.Project{}, errors.New("Error updating project: " + err.Error())
	}
	return project, nil
}

func DeleteProject(db *sql.DB, projectID int) error {
	query := `
    DELETE FROM projects WHERE id = $1;
  `
	result, err := db.Exec(query, projectID)
	if err != nil {
		return errors.New("Error deleting project: " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("Error getting rows affected: " + err.Error())
	}
	if rowsAffected == 0 {
		return errors.New("project not found")
	}

	return nil
}
