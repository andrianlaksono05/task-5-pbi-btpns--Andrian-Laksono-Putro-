package main

import (
	"log"
	"net/http"
	"pbi-final/controllers/authcontroller"
	"pbi-final/models"

	"github.com/gorilla/mux"
)

func main() {

	models.KoneksiDatabase()

	r := mux.NewRouter()

	r.HandleFunc("/users/login", authcontroller.Login).Methods("POST")
	r.HandleFunc("/users/register", authcontroller.Register).Methods("POST")
	r.HandleFunc("/users/logout", authcontroller.Logout).Methods("GET")

	r.HandleFunc("/users/photos", authcontroller.UploadPhoto).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
