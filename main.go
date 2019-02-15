package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/carsys/current-version", CurrentVersion).Methods("GET")

	myRouter.HandleFunc("/vehicle", GetVehicles).Methods("GET")
	myRouter.HandleFunc("/vehicle", CreateVehicle).Methods("POST")
	myRouter.HandleFunc("/vehicle", UpdateVehicle).Methods("PUT")
	myRouter.HandleFunc("/vehicle", DeleteVehicle).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", myRouter))
}

func main() {
	fmt.Println("Go REST API starting")

	InitialMigration()

	handleRequests()
}
