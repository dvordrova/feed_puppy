// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"database/sql"
)

type Action struct {
	ID             int64
	UserID         int64
	DogID          int64
	Timestamp      int64
	ActionID       int64
	AdditionalInfo sql.NullString
}

type ActionDescription struct {
	ID   int64
	Type string
}

type Dog struct {
	ID        int64
	Name      string
	BirthDate string
	Sex       string
	Breed     string
}

type Food struct {
	ID             int64
	DogID          int64
	Title          string
	Calories       int64
	IsLastSelected bool
}

type InviteSubscription struct {
	ID    int64
	Hash  string
	DogID int64
	Type  string
	Used  bool
}

type Subscription struct {
	ID     int64
	UserID int64
	DogID  int64
	Type   string
}

type User struct {
	ID         int64
	TelegramID int64
	Name       string
	Language   string
	CurrentDog sql.NullInt64
	State      string
}
