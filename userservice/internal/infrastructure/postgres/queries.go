package postgres

const (
	QueryFindByEmail = `
	SELECT 
		id, 
		first_name, 
		middle_name, 
		last_name, 
		hash_password, 
		email 
	FROM users 
	WHERE email = $1`
)
