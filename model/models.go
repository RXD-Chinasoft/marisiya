package model

type Friend struct {
	Id int64 `json:"id" db:"id"`
	Email string `json:"email" db:"email"`
	Friend int64 `json:"friend" db:"friend"`
}