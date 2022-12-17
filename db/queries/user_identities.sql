-- name: CreateUserIdentity :one
INSERT INTO "user_identities"(
	id,
	user_id,
	identity_hash
) VALUES (
	$1, $2, $3
) RETURNING *;

-- name: GetUserIdentityById :one
SELECT * FROM "user_identities" WHERE id = $1 LIMIT 1;

-- name: GetUserIdentityByUserId :one
SELECT * FROM "user_identities" WHERE user_id = $1 LIMIT 1;