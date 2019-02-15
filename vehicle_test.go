package main

import (
	"net/http"
	"net/url"
	"testing"
)

func Test_VehicleExists_Exists(t *testing.T) {
	vehicle := Vehicle{Licence: "11-ww-11"}
	expected := true
	actual := VehicleExists(vehicle)
	if actual != expected {
		t.Errorf("Test failed, expected: '%t', got: '%t'. Licence = %s", expected, actual, vehicle.Licence)
	}
}
func Test_VehicleExists_DoesNotExist(t *testing.T) {
	vehicle := Vehicle{}
	expected := false
	actual := VehicleExists(vehicle)
	if actual != expected {
		t.Errorf("Test failed, expected: '%t', got: '%t'. Licence = %s", expected, actual, vehicle.Licence)
	}
}

func Test_QueryParamExists_Exists(t *testing.T) {
	//Make url with parameter = 1
	u, err := url.Parse("test.com/test?parameter=1")
	if err != nil {
		panic(err.Error)
	}
	request := http.Request{URL: u}

	//Expect parameter "parameter"
	expectedString := "1"
	expectedBool := true
	actualString, actualBool := QueryParamExists("parameter", &request)

	if actualBool != expectedBool || actualString != expectedString {
		t.Errorf("Test failed, expected: '%t', '%s', got: '%t', '%s'",
			expectedBool, expectedString, actualBool, actualString)
	}
}

func Test_QueryParamExists_DoesNotExist_EmptyParamValue(t *testing.T) {
	//Make url with parameter = empty
	u, err := url.Parse("test.com/test?parameter=")
	if err != nil {
		panic(err.Error)
	}
	request := http.Request{URL: u}

	//Expect parameter "parameter"
	expectedString := ""
	expectedBool := false
	actualString, actualBool := QueryParamExists("parameter", &request)

	if actualBool != expectedBool || actualString != expectedString {
		t.Errorf("Test failed, expected: '%t', '%s', got: '%t', '%s'",
			expectedBool, expectedString, actualBool, actualString)
	}
}

func Test_QueryParamExists_DoesNotExist_NoParam1(t *testing.T) {
	//Make url without parameter
	u, err := url.Parse("test.com/test")
	if err != nil {
		panic(err.Error)
	}
	request := http.Request{URL: u}

	//Expect parameter "parameter"
	expectedString := ""
	expectedBool := false
	actualString, actualBool := QueryParamExists("parameter", &request)

	if actualBool != expectedBool || actualString != expectedString {
		t.Errorf("Test failed, expected: '%t', '%s', got: '%t', '%s'",
			expectedBool, expectedString, actualBool, actualString)
	}
}

func Test_QueryParamExists_DoesNotExist_NoParam2(t *testing.T) {
	//Make url without parameter
	u, err := url.Parse("test.com/test")
	if err != nil {
		panic(err.Error)
	}
	request := http.Request{URL: u}

	//Expect parameter ""
	expectedString := ""
	expectedBool := false
	actualString, actualBool := QueryParamExists("", &request)

	if actualBool != expectedBool || actualString != expectedString {
		t.Errorf("Test failed, expected: '%t', '%s', got: '%t', '%s'",
			expectedBool, expectedString, actualBool, actualString)
	}
}
