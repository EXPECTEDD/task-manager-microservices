package postgres

var (
	QuerieCreate            = "INSERT INTO tasks (project_id, description, deadline) VALUES($1, $2, $3) RETURNING id;"
	QuerieUpdateDescription = "UPDATE tasks SET description = $1 WHERE id = $2 AND project_id = $3;"
	QuerieUpdateDeadline    = "UPDATE tasks SET deadline = $1 WHERE id = $2 AND project_id = $3;"
)
