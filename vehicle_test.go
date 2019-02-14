package main

import "testing"

func Test_VehicleExists_Succes(t *testing.T) {
	vehicle := Vehicle{Licence: "11-ww-11"}
	expected := true
	actual := VehicleExists(vehicle)
	if actual != expected {
		t.Errorf("Test failed, expected: '%t', got: '%t'. Licence = %s", expected, actual, vehicle.Licence)
	}
}

func Test_QueryParamExists(t *testing.T) {

}
