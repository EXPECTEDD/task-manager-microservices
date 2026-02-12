package postgres

var (
	QuerieSave       = "INSERT INTO projects(owner_id, name) VALUES($1, $2) RETURNING id"
	QuerieDelete     = "DELETE FROM projects WHERE id = $1 AND owner_id = $2"
	QuerieGetAll     = "SELECT id, owner_id, name, created_at FROM projects WHERE owner_id = $1"
	QuerieUpdateName = "UPDATE projects SET name = $1 WHERE owner_id = $2 AND id = $3"
)
