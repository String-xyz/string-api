package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Geofencing interface {
	IsRestricted(ip string) (bool, error)
}

type geofencing struct {
}

func NewGeofencing() Geofencing {
	return &geofencing{}
}

type data struct {
	Ip             string `json:'ip`
	Type           string `json:'type'`
	Continent_name string `json:'continent_name'`
	Country_code   string `json:'country_code'`
	Region_code    string `json:'region_code'`
	Region_name    string `json:'region_name'`
}

func (c geofencing) IsRestricted(ip string) (bool, error) {
	fmt.Println("Client Ip Address: ", ip)

	url := "http://api.ipstack.com/" + ip + "?access_key=" + os.Getenv("LOCATION_API_KEY")

	/* ---------- Get Data from API ---------- */
	res, getErr := http.Get(url)
	if getErr != nil {
		log.Fatal(getErr)
	}

	// read the response body
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// convert the body to type string
	fmt.Println(string(body))

	data_obj := data{}

	// unmarshal the json into our struct
	jsonErr := json.Unmarshal(body, &data_obj)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	/* END ---------- Get Data from API ---------- END */

	if data_obj.Country_code != "US" || data_obj.Region_code == "NY" {
		msg := "❌ Location: " + data_obj.Region_name + " - Ip: " + ip
		fmt.Println(msg)
		return true, nil
	} else {
		msg := "✅ Location: " + data_obj.Region_name + " - Ip: " + ip
		fmt.Println(msg)
		return false, nil
	}
}
