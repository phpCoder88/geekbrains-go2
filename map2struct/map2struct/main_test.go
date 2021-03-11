package map2struct

import (
	"errors"
	"reflect"
	"testing"
)

type User struct {
	UserName string
	UserAge  int
	IsAdmin  bool
	Roles    []string
	Test     map[string]int
	Address  Address
}

type Address struct {
	Country   string
	City      string
	Street    string
	Building  int
	Apartment int
	Location  Location
}

type Location struct {
	Latitude  float64
	Longitude float64
}

var (
	inStruct = User{
		UserName: "Test",
		UserAge:  33,
		IsAdmin:  false,
		Roles:    []string{"Guest"},
		Test:     map[string]int{"123": 123, "547": 547},
	}

	inValues = map[string]interface{}{
		"UserName": "Pavel",
		"IsAdmin":  true,
		"Roles":    []string{"Admin"},
		"Test":     map[string]int{"345": 345, "890": 890},
		"Address": Address{
			Country:   "Russian",
			City:      "Moscow",
			Street:    "Kremlin st",
			Building:  24,
			Apartment: 145,
			Location: Location{
				Latitude:  131234545.4358,
				Longitude: 657567323.5667,
			},
		},
	}

	wantOutStruct = User{
		UserName: "Pavel",
		UserAge:  33,
		IsAdmin:  true,
		Roles:    []string{"Admin"},
		Test:     map[string]int{"345": 345, "890": 890},
		Address: Address{
			Country:   "Russian",
			City:      "Moscow",
			Street:    "Kremlin st",
			Building:  24,
			Apartment: 145,
			Location: Location{
				Latitude:  131234545.4358,
				Longitude: 657567323.5667,
			},
		},
	}
)

func TestMapToStruct_InnerStructs(t *testing.T) {
	err := MapToStruct(&inStruct, inValues)

	if err != nil {
		t.Errorf("test Failed - results not match\nGot:\n%v\nExpected:\n%v", err, nil)
	}

	if !reflect.DeepEqual(inStruct, wantOutStruct) {
		t.Errorf("test Failed - results not match\nGot:\n%v\nExpected:\n%v", inStruct, wantOutStruct)
	}

}

func TestMapToStruct_NotPointerError(t *testing.T) {
	err := MapToStruct(inStruct, inValues)

	if err != ErrNotPinter {
		t.Errorf("test Failed - results not match\nGot:\n%v\nExpected:\n%v", err, ErrNotPinter)
	}
}

func TestMapToStruct_NotStruct(t *testing.T) {

	inNotStruct := 123
	err := MapToStruct(&inNotStruct, inValues)

	if err != ErrNotStruct {
		t.Errorf("test Failed - results not match\nGot:\n%v\nExpected:\n%v", err, ErrNotStruct)
	}
}

func TestMapToStruct_IncompatibleType(t *testing.T) {
	inStruct := struct {
		Name string
		Age  int
	}{
		Name: "Pavel",
		Age:  33,
	}
	inValues := map[string]interface{}{
		"Name": "Pavel",
		"Age":  "33",
	}
	err := MapToStruct(&inStruct, inValues)

	if !errors.Is(err, ErrIncompatibleType) {
		t.Errorf("test Failed - results not match\nGot:\n%v\nExpected:\n%v", err, ErrIncompatibleType)
	}
}

func TestMapToStruct_CannotSetValue(t *testing.T) {
	inStruct := struct {
		Name string
		age  int
	}{
		Name: "Pavel",
		age:  33,
	}
	inValues := map[string]interface{}{
		"Name": "Pavel",
		"age":  33,
	}
	err := MapToStruct(&inStruct, inValues)

	if !errors.Is(err, ErrCannotSetValue) {
		t.Errorf("test Failed - results not match\nGot:\n%v\nExpected:\n%v", err, ErrCannotSetValue)
	}
}
