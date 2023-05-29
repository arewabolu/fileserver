package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func groupRoutes(r chi.Router) {
	r.Use(middlewareAUth)

	r.Get("/user/list-files/", UserPageHandle)
	r.Get("/user/download/", UserDownloadHandle)

	r.Post("/user/upload/{path}/", UserUploadHandle)
	r.Post("/user/create-new-folder/", createFolderHandle)
}

func main() {
	router := chi.NewRouter()

	router.Use(middleware.CleanPath)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Group(groupRoutes)
	router.Post("/register/", CreateAccount)
	router.Post("/login/", authHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}

// Aim is to certify that the person sending a url
// is the same one in the cookie id
