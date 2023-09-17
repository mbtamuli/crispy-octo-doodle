-- name: GetClient :one
SELECT * FROM clients
WHERE id = ? LIMIT 1;

-- name: ListClients :many
SELECT * FROM clients
ORDER BY name;

-- name: CreateClient :execresult
INSERT INTO clients (
  name, email
) VALUES (
  ?, ?
);

-- name: UpdateClient :execresult
UPDATE clients
SET plan_id = $2
WHERE id = $1;

-- name: GetKeys :one
SELECT * FROM authkeys
WHERE client_id = ? LIMIT 1;

-- name: CreateKey :execresult
INSERT INTO authkeys (
  client_id, access, secret
) VALUES (
  ?, ?, ?
);

-- name: GetImage :one
SELECT * FROM images
WHERE id = ? LIMIT 1;

-- name: ListImages :many
SELECT * FROM images
ORDER BY path;

-- name: ListImagesForClient :many
SELECT * FROM images
WHERE client_id = $1;

-- name: CreateImage :execresult
INSERT INTO images (
  path, size, type, extension, client_id
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: CreateMatch :execresult
INSERT INTO matches (
  image1_id, image2_id, score, client_id
) VALUES (
  ?, ?, ?, ?
);

-- name: ListMatchesForClient :many
SELECT * FROM matches
WHERE client_id = $1;