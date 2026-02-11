package postgres

var (
	QuerieCreate = "INSERT INTO tasks (project_id, description, deadline) VALUES($1, $2, $3) RETURNING id;"
)
