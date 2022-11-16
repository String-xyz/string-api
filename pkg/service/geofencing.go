package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/String-xyz/string-api/pkg/store"
)

const A_DAY_IN_NANOSEC = 84600000000000

type Geofencing interface {
	IsAllowed(ip string) (bool, error)
	getLocationDataFromRedis(ip string) (locationData, error)
	getLocationDataFromAPI(ip string) (locationData, error)
	setLocationDataInRedis(ip string, locationData locationData) error
}

type geofencing struct {
	redis store.RedisStore
}

func NewGeofencing(redis store.RedisStore) Geofencing {
	return &geofencing{redis}
}

type locationData struct {
	Ip           string `json:'ip`
	Country_code string `json:'country_code'`
	Region_code  string `json:'region_code'`
	Region_name  string `json:'region_name'`
}

func (c geofencing) IsAllowed(ip string) (bool, error) {
	fmt.Println("Client Ip Address: ", ip)
	locationData, err := c.getLocationDataFromRedis(ip)

	/* if not found in cache, get it from the api then set the value to api's response*/
	if err != nil {
		locationData, err = c.getLocationDataFromAPI(ip)
		if err != nil {
			// TODO: log error
			return false, err
		}

		err = c.setLocationDataInRedis(ip, locationData)
		if err != nil { // TODO: log error
			// TODO: log error
			return false, err
		}
	}

	isAllowed := isRegionAllowed(locationData)
	if isAllowed {
		fmt.Println("✅ Location: " + locationData.Region_name + " - Ip: " + ip)

	} else {
		fmt.Println("❌ Location: " + locationData.Region_name + " - Ip: " + ip)
	}

	return isAllowed, nil
}

func (c geofencing) setLocationDataInRedis(ip string, locationData locationData) error {
	locationDataStr, err := json.Marshal(locationData)
	if err != nil {
		return err
	}

	return c.redis.Set("ip:"+ip, locationDataStr, A_DAY_IN_NANOSEC)
}

func (c geofencing) getLocationDataFromRedis(ip string) (locationData, error) {
	cachedData, err := c.redis.Get("ip:" + ip)
	locationData := locationData{}

	if cachedData == nil {
		return locationData, err
	}
	err = json.Unmarshal(cachedData, &locationData)
	if err != nil {
		return locationData, err
	}

	if locationData.Ip != ip {
		return locationData, fmt.Errorf("ip not found in cache")
	}

	return locationData, err
}

func (c geofencing) getLocationDataFromAPI(ip string) (locationData, error) {
	url := "http://api.ipstack.com/" + ip + "?access_key=" + os.Getenv("LOCATION_API_KEY")

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

	data_obj := locationData{}

	// unmarshal the json into our struct
	jsonErr := json.Unmarshal(body, &data_obj)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return data_obj, nil
}

func isRegionAllowed(locationData locationData) bool {
	return locationData.Country_code == "US" && locationData.Region_code != "NY"
}
