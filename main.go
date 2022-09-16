package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func groupRoutes(r chi.Router) {
	r.Use(middlewareAUth)
	r.Get("/user/{ID}/", UserPageHandle)
	r.Get("/user/{ID}/download/", UserDownloadHandle)
	r.Post("/user/{ID}/upload/", UserUploadHandle)
	r.Post("/user/{ID}/create-folder/", createFolderHandle)
}

func main() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Group(groupRoutes)
	router.Post("/register/", postCreateAccount)
	router.Post("/login/", authHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}

// Aim is to certify that the person sending a url
// is the same one in the cookie id
