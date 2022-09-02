package main

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/go-pg/pg/v10"
)

func homeHandle(w http.ResponseWriter, r *http.Request) {

}

func loginHandle(w http.ResponseWriter, r *http.Request) {

}

func postCreateAccount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	hashedPassword, err := passwordHash(r.PostFormValue("password"))
	if err != nil {
		fmt.Fprint(w, err)
	}
	creds := &User{
		Email:        strings.ToLower(r.PostFormValue("email")),
		PasswordHash: string(hashedPassword),
		FullName:     r.PostFormValue("fname"),
	}
	fmt.Println(creds.Email)
	queryErr := dbConn.Model(creds).Column("email").Where("?=?", pg.Ident("email"), creds.Email).Select()
	if queryErr != nil {
		fmt.Println("trying to create account")
		fmt.Println(queryErr)
		//convert names and email to lower case for db
		_, insErr := dbConn.Model(creds).Insert()
		//	err := dbMan.QueryRow(insertStatement, fullName, ), email)
		if insErr != nil {
			fmt.Println("no error before here")
			fmt.Println(insErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "localhost:8080/login/", http.StatusPermanentRedirect)
		return

	}
	w.WriteHeader(http.StatusInternalServerError)
}
func createAccount(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("./assets/createAcc.html")
	t.Execute(w, nil)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	email := strings.ToLower(r.PostFormValue("email"))
	Password := r.PostFormValue("password")
	if email == "" {
		w.WriteHeader(http.StatusExpectationFailed)
	}
	if Password == "" {
		w.WriteHeader(http.StatusExpectationFailed)
	}

	dbCred := &User{}
	//scan into dbcred from database
	SelectErr := dbConn.Model(dbCred).Column("passwordhash", "id").Where("email=?", email).Select(dbCred.PasswordHash, dbCred.ID)
	if SelectErr != nil {
		w.Write([]byte(SelectErr.Error()))
	}

	pswdValErr := validateHash(dbCred.PasswordHash, Password)
	if pswdValErr != nil {
		w.Write([]byte(pswdValErr.Error()))
	}

	JWTtoken, err := GenerateJWT(email)
	if err != nil {
		w.Write([]byte("User does not exist"))
		http.Redirect(w, r, "localhost:8080/reg/", http.StatusPermanentRedirect)
	}

	expirationTime := time.Now().Add(2 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Value:   JWTtoken,
		Expires: expirationTime,
	})
	fmt.Println(dbCred.ID, dbCred.FullName)
	w.Header().Add("jwt", JWTtoken)

	//http.Redirect(w, r, "localhost:8080/user/{dbCreds.ID}", http.StatusPermanentRedirect)
}

func UserPageHandle(w http.ResponseWriter, r *http.Request) {

}
