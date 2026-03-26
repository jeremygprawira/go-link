package pgsql

var (
	// QueryCreateUrl inserts request to the urls table
	QueryCreateUrl = `
		INSERT into urls (id, code, name, url, account_number, click_count, state, metadata, expired_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now(), now())
	`

	// QueryCheckByCode checks if a code already exists in the urls table
	QueryCheckByCode = `
		SELECT EXISTS(SELECT 1 FROM urls WHERE code = $1)
	`

	// QueryGetOneByCode gets a user by code (shortened code).
	QueryGetOneByCode = `
		SELECT id, code, name, url, account_number, click_count, state, metadata, expired_at, created_at, updated_at FROM users
		WHERE code = $1 AND deleted_at IS NULL
	`

	QueryGetManyByAccountNumber = `
		SELECT id, code, name, url, account_number, click_count, state, metadata, expired_at, created_at, updated_at FROM users
		WHERE code = $1 AND deleted_at IS NULL
	`

	// QueryGetOneByCode gets a user by code (shortened code).
	QueryGetOneById = `
		SELECT id, code, name, url, account_number, click_count, state, metadata, expired_at, created_at, updated_at FROM users
		WHERE code = $1 AND deleted_at IS NULL
	`

	// QueryUpdateMetadataByCode update link metadata by code (shortened code)
	QueryUpdateMetadataByCode = `
		UPDATE urls
		SET metadata = $1
		WHERE code = $2
		RETURNING metadata
	`

	// QueryUpdateClickCountByCode update click count (by increment) by code (shortened code)
	QueryUpdateClickCountByCode = `
		UPDATE urls
		SET click_count = $1
		WHERE code = $2
		RETURNING click_count
	`

	// QueryUpdateMetadataByCode update link metadata by id (uuid)
	QueryUpdateMetadataById = `
		UPDATE urls
		SET metadata = $1
		WHERE id = $2
		RETURNING metadata
	`

	// QueryUpdateClickCountByCode update click count (by increment) by id (uuid)
	QueryUpdateClickCountById = `
		UPDATE urls
		SET click_count = $1
		WHERE id = $2
		RETURNING click_count
	`

	// QueryDeleteById
	QueryDeleteById = `
		UPDATE urls
		SET deleted_at = NOW()
		WHERE id = $1
	`
)
