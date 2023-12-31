// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package database

import (
	"database/sql"
)

type Authkey struct {
	ClientID int32
	Access   string
	Secret   string
}

type Client struct {
	ID     int32
	Name   string
	Email  string
	PlanID sql.NullInt32
}

type Image struct {
	ID        int32
	Path      string
	Size      string
	Type      string
	Extension string
	ClientID  sql.NullInt32
}

type Match struct {
	Image1ID int32
	Image2ID int32
	Score    int32
	ClientID sql.NullInt32
}

type Plan struct {
	ID              int32
	Name            string
	Base            int32
	FaceMatchAndOcr int32
	Upload          int32
}
