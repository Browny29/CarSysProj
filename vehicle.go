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
	db, err = gorm.Open(databaseType, connectionString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Could not connect to database")
	}

	defer db.Close()

	db.AutoMigrate(&Vehicle{})
}

func GetVehicles(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open(databaseType, connectionString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Could not connect to database")
	}

	defer db.Close()

	//Check the parameter form the URI
	licenceplate, ParamExists := DoesQueryParamExist("licenceplate", r)
	if !ParamExists {
		//GET all vehicles
		var vehicles []Vehicle
		db.Find(&vehicles)
		json.NewEncoder(w).Encode(vehicles)
	} else {
		//GET specific vehicle
		var vehicle Vehicle
		db.Where("licence = ?", licenceplate).Find(&vehicle)

		//Check if vehicle exists
		if vehicle.ID != 0 {
			//Show vehicle
			json.NewEncoder(w).Encode(vehicle)
		} else {
			//Give error message back
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Vehicle not found")
		}

	}

}

func CreateVehicle(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open(databaseType, connectionString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Could not connect to database")
	}

	defer db.Close()

	var vehicle Vehicle

	DecodeJsonObject(r.Body, &vehicle)

	//Create new record in the database
	output := db.Create(&Vehicle{Licence: vehicle.Licence, Brand: vehicle.Brand, Type: vehicle.Type, Buildyear: vehicle.Buildyear,
		Odometer: vehicle.Odometer, Unit: vehicle.Unit, Color: vehicle.Color, Weight: vehicle.Weight})

	//Convert the vehicle back to json and write it in the response
	vehicleJson, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	//Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(vehicleJson)
}

func UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open(databaseType, connectionString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Could not connect to database")
	}

	defer db.Close()
	//Check if URI parameter is correct
	licenceplate, ParamExists := DoesQueryParamExist("licenceplate", r)
	if ParamExists {
		//Update Vehicle
		var vehicle Vehicle
		var newVehicle Vehicle

		db.Where("licence = ?", licenceplate).Find(&vehicle)

		//Check if vehicle exists
		if vehicle.ID != 0 {
			DecodeJsonObject(r.Body, &newVehicle)

			//Update values of vehicle
			newVehicle.ID = vehicle.ID
			newVehicle.CreatedAt = vehicle.CreatedAt
			output := db.Save(&newVehicle)

			//Convert the vehicle back to json and write it in the response
			vehicleJson, err := json.Marshal(output)
			if err != nil {
				panic(err)
			}

			//Write the response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(vehicleJson)

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
	db, err = gorm.Open(databaseType, connectionString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Could not connect to database")
	}

	defer db.Close()

	//Check the parameter form the URI
	licenceplate, ParamExists := DoesQueryParamExist("licenceplate", r)
	if ParamExists {
		//Delete Vehicle
		var vehicle Vehicle
		db.Where("licence = ?", licenceplate).Find(&vehicle)

		//Check if vehicle exists
		if vehicle.ID != 0 {
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

func OpenDatabaseConnection(db *gorm.DB) *gorm.DB {
	db, err = gorm.Open(databaseType, connectionString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Could not connect to database")
	}

	defer db.Close()

	return db
}

func CloseDatabaseConnection(db *gorm.DB) {
	defer db.Close()
}

func DecodeJsonObject(jsonInput io.ReadCloser, vehicle *Vehicle) {
	err := json.NewDecoder(jsonInput).Decode(&vehicle)
	if err != nil {
		panic(err)
	}
}

func DoesQueryParamExist(parameter string, r *http.Request) (string, bool) {
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
