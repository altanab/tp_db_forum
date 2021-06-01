package models

import (
	"database/sql"
	"encoding/json"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type UserIDError struct {
	Message string `json:"message"`
}

func (e UserIDError) Error() string {
	return e.Message
}

type AlreadyExists struct {
	Message string `json:"message"`
}

func (e AlreadyExists) Error() string {
	return e.Message
}

type NullInt struct {
	sql.NullInt64
}

type NullString struct {
	sql.NullString
}

type User struct {
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About string `json:"about"`
	Email string `json:"email"`
}

type Forum struct {
	Title string `json:"title"`
	User string `json:"user"`
	Slug string `json:"slug"`
	Posts int `json:"posts"`
	Threads int `json:"threads"`
}

type Thread struct {
	Author  string         `json:"author"`
	Created time.Time      `json:"created"`
	Forum   string         `json:"forum"`
	Id      int            `json:"id"`
	Message string         `json:"message"`
	Slug    NullString      `json:"slug"`
	Title   string         `json:"title"`
	Votes   int            `json:"votes"`
}


type ThreadUpdate struct {
	Title string `json:"title"`
	Message string `json:"message"`
}

type Post struct {
	Id int `json:"id"`
	Parent NullInt `json:"parent"`
	Author string `json:"author"`
	Message string `json:"message"`
	IsEdited bool `json:"isEdited"`
	Forum string `json:"forum"`
	Thread int `json:"thread"`
	Created time.Time `json:"created"`
}


type PostUpdate struct {
	Id int `-"`
	Message string `json:"message"`
}

type PostFull struct {
	Post Post
	Author User
	Thread Thread
	Forum Forum
}

type Vote struct {
	Nickname string `json:"nickname"`
	Voice int `json:"voice"`
}

type Status struct {
	User int `json:"user"`
	Forum int `json:"forum"`
	Thread int `json:"thread"`
	Post int `json:"post"`
}

var DBConn *pgxpool.Pool

func (ns NullInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(ns.Int64)
}

func (ns *NullInt) UnmarshalJSON(data []byte) error {
	var b *int64
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if *b != 0 {
		ns.Valid = true
		ns.Int64 = *b
	} else {
		ns.Valid = false
	}
	return nil
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var b *string
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if b != nil {
		ns.Valid = true
		ns.String = *b
	} else {
		ns.Valid = false
	}
	return nil
}