package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Version struct {
	Build_date string `json:"Build date"`
}

func currentVersion(w http.ResponseWriter, r *http.Request) {
	var version Version
	resp, err := http.Get("http://nl.carsys.online/version.json")
	if err != nil {
		panic(err.Error)
	}

	json.NewDecoder(resp.Body).Decode(&version)
	fmt.Fprintf(w, version.Build_date)
}
