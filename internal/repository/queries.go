package repository

const (
	selectAll = `
		SELECT * FROM orders;
	`
	createOrderQuery = `
		INSERT INTO
			orders
			(details, created_at, updated_at)
		VALUES
			($1, $2, $3)
		RETURNING
			id;
	`
	findByID = `
		SELECT
			id, details, created_at, updated_at
		FROM
			orders
		WHERE
			id = $1;
	`
	updateByID = `
		UPDATE
			orders
		SET
			details = $1,
			updated_at = $2
		WHERE 
			id = $3; 
	`
	deleteByID = `
		DELETE FROM
			orders
		WHERE
			id = $1;
	`
)
