package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-pg/pg/v10"
)

const maxUploadSize = 200 << 20

//const maxUploadSize2 = 20 * 10 * 1024 * 1024

func postCreateAccount(w http.ResponseWriter, r *http.Request) {
	creds := &User{}
	jsErr := json.NewDecoder(r.Body).Decode(creds)
	if jsErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(creds)

	queryErr := dbConn.Model(creds).Column("email").Where("?=?", pg.Ident("email"), strings.ToLower(creds.Email)).Select()
	if queryErr != nil {
		//convert names and email to lower case for db
		hashedPaswd, _ := passwordHash(creds.PasswordHash)
		dbCred := &User{
			Email:        strings.ToLower(creds.Email),
			FullName:     strings.ToLower(creds.FullName),
			PasswordHash: string(hashedPaswd),
		}

		//Miracle in Cell No. 7, The Grand Heist

		_, insErr := dbConn.Model(dbCred).Insert()
		if insErr != nil {
			w.Write([]byte("unfortunately we couldn't create your account. Please try again."))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		SelectErr := dbConn.Model(dbCred).Column("id").Where("email=?", dbCred.Email).Select(&dbCred.ID)
		if SelectErr != nil {
			fmt.Println("post create error:", SelectErr)
			return
		}

		err := createUserBucket(strAppender(strconv.Itoa(dbCred.ID)))
		if err != nil {
			fmt.Println("create bucket error:", err)
		}
		//	http.Redirect(w, r, "localhost:8080/login", http.StatusPermanentRedirect)
		return
	}

	w.Write([]byte("account already exists. please login"))
	w.WriteHeader(http.StatusOK)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	creds := &User{}
	json.NewDecoder(r.Body).Decode(creds)

	email := strings.ToLower(creds.Email)
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	Password := creds.PasswordHash
	if Password == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	dbCred := &User{}
	//scan into dbcred from database
	SelectErr := dbConn.Model(dbCred).Column("passwordhash", "id").Where("email=?", email).Select(&dbCred.PasswordHash, &dbCred.ID)
	if SelectErr != nil {
		w.Write([]byte("please enter a valid email"))
	}

	pswdValErr := validateHash(Password, dbCred.PasswordHash)
	if pswdValErr != nil {
		w.Write([]byte("incorrect password"))
	}
	JWTtoken, err := GenerateJWT(dbCred.ID)
	if err != nil {
		w.Write([]byte("user does not exist")) //why here?
	}

	w.Header().Add("Authorization", "Bearer"+JWTtoken)
	s := fmt.Sprintln("your user ID is "+strconv.Itoa(dbCred.ID), "please use this for subsequent request to upload or download")
	w.Write([]byte(s))

	//	fmt.Println(dbCred.ID, dbCred.FullName)
}

func UserPageHandle(w http.ResponseWriter, r *http.Request) {
	//get data from S3 to display as json response

	list, err := listFiles(reqSplit(r))
	if err != nil {
		w.Write([]byte("failed to load list"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jL, err := json.Marshal(list)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(jL)
	w.WriteHeader(http.StatusOK)
}

func UserDownloadHandle(w http.ResponseWriter, r *http.Request) {
	//get from s3 and return to user
	w.Header().Set("Content-Type", "application/octet-stream")
	file := new(Files)
	json.NewDecoder(r.Body).Decode(file)
	keyName := reqSplit(r) + "/" + file.Filename
	usrFile := downloadFile(keyName)

	byteFile, _ := io.ReadAll(usrFile)
	w.Write(byteFile)
}

func UserUploadHandle(w http.ResponseWriter, r *http.Request) {
	//parse data and send to s3 bucket

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	err := r.ParseMultipartForm(maxUploadSize)

	if err != nil {
		if err == http.ErrNotMultipart {
			w.Write([]byte("no file uploaded. Please confirm your upload"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Write([]byte("the file uploaded is too large. Please select a file not larger than 200MB!"))
		return
	}

	file, fileHeader, err := r.FormFile("upload-file")
	if err == http.ErrMissingFile {
		w.Write([]byte("unable to upload file. Please check upload info"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	err = uploadFile(reqSplit(r)+"/"+fileHeader.Filename, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write([]byte("yourfile has been sucessfully uploaded"))
	w.WriteHeader(http.StatusOK)
}

func createFolderHandle(w http.ResponseWriter, r *http.Request) {
	var folderName string
	err := json.NewDecoder(r.Body).Decode(&folderName)
	if err != nil {
		w.Write([]byte("could not process your request. please try again."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = createFolder(folderName)
	if err != nil {
		w.Write([]byte("we were unable to create your folder"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
