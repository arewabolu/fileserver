package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", homeHandle)

	router.Get("/login/", loginHandle)

	router.Get("/user/*", UserPageHandle)
	router.Get("/reg/", createAccount)
	router.Post("/reg/", postCreateAccount)
	router.Post("/login/", authHandler)

	//homeHandler := http.FileServer(http.Dir("./assets/root"))
	//router.Handle("/", http.StripPrefix("/", homeHandler))
	//router.HandleFunc("/", homeHandle)
	//router.HandleFunc()
	//router.HandleFunc()
	//	router.HandleFunc("", homeHandle)

	log.Fatal(http.ListenAndServe(":8080", router))
}
