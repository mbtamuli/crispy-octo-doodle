// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: query.sql

package database

import (
	"context"
	"database/sql"
)

const createClient = `-- name: CreateClient :execlastid
INSERT INTO clients (
  name, email, plan_id
) VALUES (
  ?, ?, ?
)
`

type CreateClientParams struct {
	Name   string
	Email  string
	PlanID sql.NullInt32
}

func (q *Queries) CreateClient(ctx context.Context, arg CreateClientParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, createClient, arg.Name, arg.Email, arg.PlanID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

const createImage = `-- name: CreateImage :execlastid
INSERT INTO images (
  path, size, type, extension, client_id
) VALUES (
  ?, ?, ?, ?, ?
)
`

type CreateImageParams struct {
	Path      string
	Size      string
	Type      string
	Extension string
	ClientID  sql.NullInt32
}

func (q *Queries) CreateImage(ctx context.Context, arg CreateImageParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, createImage,
		arg.Path,
		arg.Size,
		arg.Type,
		arg.Extension,
		arg.ClientID,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

const createKey = `-- name: CreateKey :execresult
INSERT INTO authkeys (
  client_id, access, secret
) VALUES (
  ?, ?, ?
)
`

type CreateKeyParams struct {
	ClientID int32
	Access   string
	Secret   string
}

func (q *Queries) CreateKey(ctx context.Context, arg CreateKeyParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createKey, arg.ClientID, arg.Access, arg.Secret)
}

const createMatch = `-- name: CreateMatch :execresult
INSERT INTO matches (
  image1_id, image2_id, score, client_id
) VALUES (
  ?, ?, ?, ?
)
`

type CreateMatchParams struct {
	Image1ID int32
	Image2ID int32
	Score    int32
	ClientID sql.NullInt32
}

func (q *Queries) CreateMatch(ctx context.Context, arg CreateMatchParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createMatch,
		arg.Image1ID,
		arg.Image2ID,
		arg.Score,
		arg.ClientID,
	)
}

const getClientID = `-- name: GetClientID :one
SELECT id FROM clients
WHERE email = ? LIMIT 1
`

func (q *Queries) GetClientID(ctx context.Context, email string) (int32, error) {
	row := q.db.QueryRowContext(ctx, getClientID, email)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const getImage = `-- name: GetImage :one
SELECT id, path, size, type, extension, client_id FROM images
WHERE id = ? LIMIT 1
`

func (q *Queries) GetImage(ctx context.Context, id int32) (Image, error) {
	row := q.db.QueryRowContext(ctx, getImage, id)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.Path,
		&i.Size,
		&i.Type,
		&i.Extension,
		&i.ClientID,
	)
	return i, err
}

const getKeys = `-- name: GetKeys :one
SELECT client_id, access, secret FROM authkeys
WHERE client_id = ? LIMIT 1
`

func (q *Queries) GetKeys(ctx context.Context, clientID int32) (Authkey, error) {
	row := q.db.QueryRowContext(ctx, getKeys, clientID)
	var i Authkey
	err := row.Scan(&i.ClientID, &i.Access, &i.Secret)
	return i, err
}

const listClients = `-- name: ListClients :many
SELECT id, name, email, plan_id FROM clients
ORDER BY name
`

func (q *Queries) ListClients(ctx context.Context) ([]Client, error) {
	rows, err := q.db.QueryContext(ctx, listClients)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Client
	for rows.Next() {
		var i Client
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.PlanID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listImages = `-- name: ListImages :many
SELECT id, path, size, type, extension, client_id FROM images
ORDER BY path
`

func (q *Queries) ListImages(ctx context.Context) ([]Image, error) {
	rows, err := q.db.QueryContext(ctx, listImages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Image
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.ID,
			&i.Path,
			&i.Size,
			&i.Type,
			&i.Extension,
			&i.ClientID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listImagesForClient = `-- name: ListImagesForClient :many
SELECT id, path, size, type, extension, client_id FROM images
WHERE client_id = $1
`

func (q *Queries) ListImagesForClient(ctx context.Context) ([]Image, error) {
	rows, err := q.db.QueryContext(ctx, listImagesForClient)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Image
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.ID,
			&i.Path,
			&i.Size,
			&i.Type,
			&i.Extension,
			&i.ClientID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listMatchesForClient = `-- name: ListMatchesForClient :many
SELECT image1_id, image2_id, score, client_id FROM matches
WHERE client_id = $1
`

func (q *Queries) ListMatchesForClient(ctx context.Context) ([]Match, error) {
	rows, err := q.db.QueryContext(ctx, listMatchesForClient)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Match
	for rows.Next() {
		var i Match
		if err := rows.Scan(
			&i.Image1ID,
			&i.Image2ID,
			&i.Score,
			&i.ClientID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPlans = `-- name: ListPlans :many
SELECT id, name, base, face_match_and_ocr, upload FROM plans
`

func (q *Queries) ListPlans(ctx context.Context) ([]Plan, error) {
	rows, err := q.db.QueryContext(ctx, listPlans)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Plan
	for rows.Next() {
		var i Plan
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Base,
			&i.FaceMatchAndOcr,
			&i.Upload,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateClient = `-- name: UpdateClient :execresult
UPDATE clients
SET plan_id = $2
WHERE id = $1
`

func (q *Queries) UpdateClient(ctx context.Context) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateClient)
}
