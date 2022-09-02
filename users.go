package main

type User struct {
	ID           int    `json:"id"`
	FullName     string `json: "fullname"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordhash"`
}

//type Files struct {
//	ID        int      `pg:`
//	FileStore []string `pg:`
//}
