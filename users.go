package main

type User struct {
	ID           int    `json:"-"`
	FullName     string `json:"fullname" pg:"fullname"`
	Email        string `json:"email"`
	PasswordHash string `json:"password" pg:"passwordhash"`
}

type Files struct {
	Filename string `json:"filename"`
}
