package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error
var connectionString = "host=localhost port=5432 user=postgres dbname=carsys password=root sslmode=disable"
var databaseType = "postgres"

type Vehicle struct {
	gorm.Model
	Licence   string
	Brand     string
	Type      string
	Buildyear int
	Odometer  int
	Unit      string
	Color     string
	Weight    int
}

func InitialMigration() {
	OpenDatabaseConnection()
	//defer db.Close()

	db.AutoMigrate(&Vehicle{})
}

func GetVehicles(w http.ResponseWriter, r *http.Request) {
	//OpenDatabaseConnection()
	//defer db.Close()

	//Check the parameter form the URI
	licenceplate, ParamExists := QueryParamExist("licenceplate", r)
	if !ParamExists {
		vehicles := EncodeJsonObject(GetAllVehicles(w))
		WriteSuccessResponse(w, vehicles, "json")
	} else {
		vehicle := GetSpecificVehicle(w, licenceplate)
		jsonVehicle := EncodeJsonObject(vehicle)
		if VehicleExists(vehicle) {
			WriteSuccessResponse(w, jsonVehicle, "json")
		} else {
			WriteNotFoundResponse(w)
		}
	}
}
func GetAllVehicles(w http.ResponseWriter) []Vehicle {
	var vehicles []Vehicle
	db.Find(&vehicles)
	return vehicles
}
func GetSpecificVehicle(w http.ResponseWriter, licenceplate string) Vehicle {
	var vehicle Vehicle
	db.Where("licence = ?", licenceplate).Find(&vehicle)

	return vehicle
}

func CreateVehicle(w http.ResponseWriter, r *http.Request) {
	//OpenDatabaseConnection()
	//defer db.Close()

	var vehicle Vehicle

	DecodeJsonObject(r.Body, &vehicle)

	//Create new record in the database
	output := db.Create(&Vehicle{Licence: vehicle.Licence, Brand: vehicle.Brand, Type: vehicle.Type, Buildyear: vehicle.Buildyear,
		Odometer: vehicle.Odometer, Unit: vehicle.Unit, Color: vehicle.Color, Weight: vehicle.Weight})

	//Convert the vehicle back to json and write it in the response
	vehicleJson := EncodeJsonObject(output)
	WriteSuccessResponse(w, vehicleJson, "json")
}

func UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	//OpenDatabaseConnection()
	//defer db.Close()

	//Check if URI parameter is correct
	licenceplate, ParamExists := QueryParamExist("licenceplate", r)
	if ParamExists {
		//Update Vehicle
		var vehicle Vehicle
		var newVehicle Vehicle

		db.Where("licence = ?", licenceplate).Find(&vehicle)

		if VehicleExists(vehicle) {
			DecodeJsonObject(r.Body, &newVehicle)

			//Update values of vehicle
			newVehicle.ID = vehicle.ID
			newVehicle.CreatedAt = vehicle.CreatedAt
			output := db.Save(&newVehicle)

			//Convert the vehicle back to json and write it in the response
			vehicleJson := EncodeJsonObject(output)
			WriteSuccessResponse(w, vehicleJson, "json")

		} else {
			//Give error message back
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Vehicle not found. Nothing had been updated")
		}
	} else {
		//Faulty parameter format
		w.WriteHeader(http.StatusBadRequest)
	}
}

func DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	//OpenDatabaseConnection()
	//defer db.Close()

	//Check the parameter form the URI
	licenceplate, ParamExists := QueryParamExist("licenceplate", r)
	if ParamExists {
		//Delete Vehicle
		var vehicle Vehicle
		db.Where("licence = ?", licenceplate).Find(&vehicle)

		if VehicleExists(vehicle) {
			//Delete vehicle
			db.Delete(&vehicle)
			fmt.Fprintf(w, "Vehicle with licenceplate "+licenceplate+" deleted")
		} else {
			//Give error message back
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Vehicle not found. Nothing had been deleted")
		}
	} else {
		//Fault parameter format
		w.WriteHeader(http.StatusBadRequest)
	}
}

func WriteSuccessResponse(w http.ResponseWriter, body []byte, contentType string) {
	switch contentType {
	case "json":
		w.Header().Set("Content-Type", "application/json")
	default:
		panic("Content-type not recognized")
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func WriteNotFoundResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Vehicle not found")
}

func OpenDatabaseConnection() {
	db, err = gorm.Open(databaseType, connectionString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Could not connect to database")
	}
}

func CloseDatabaseConnection(db *gorm.DB) {
	defer db.Close()
}

func EncodeJsonObject(object interface{}) []byte {
	jsonObject, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}
	return jsonObject
}

func DecodeJsonObject(jsonInput io.ReadCloser, vehicle *Vehicle) {
	err := json.NewDecoder(jsonInput).Decode(&vehicle)
	if err != nil {
		panic(err)
	}
}

func VehicleExists(vehicle Vehicle) bool {
	if vehicle.ID != 0 {
		return true
	} else {
		return false
	}
}

func QueryParamExist(parameter string, r *http.Request) (string, bool) {
	keys, ok := r.URL.Query()[parameter]

	if !ok || len(keys[0]) < 1 {
		//The parameter does not exist
		return "", false
	} else {
		//Return the value of the parameter
		key := keys[0]
		return string(key), true
	}
}
