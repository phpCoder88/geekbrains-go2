package main

import (
	"fmt"

	"github.com/phpCoder88/geekbrains-go2/map2struct/map2struct"
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

func main() {

	structIn := User{
		UserName: "Test",
		UserAge:  5,
		IsAdmin:  false,
		Roles:    []string{"Guest"},
		Test:     map[string]int{"123": 123, "547": 547},
	}

	values := map[string]interface{}{
		"UserName": "Pavel",
		"UserAge":  33,
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

	fmt.Println(map2struct.MapToStruct(&structIn, values))
	fmt.Printf("%#v\n", structIn)
}
