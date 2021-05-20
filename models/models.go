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
	Id int `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Forum string `json:"forum"`
	Message string `json:"message"`
	Votes int `json:"votes"`
	Slug string `json:"slug"`
	Created time.Time `json:"created"`
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

func (ns *NullInt) MarshalJSON() ([]byte, error) {
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