package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/pkg/errors"
)

const A_DAY_IN_NANOSEC = 84600000000000

type Geofencing interface {
	IsAllowed(ip string) (bool, error)
	getLocationDataFromRedis(ip string) (locationData, error)
	setLocationDataInRedis(ip string, locationData locationData) error
}

type geofencing struct {
	redis store.RedisStore
}

func NewGeofencing(redis store.RedisStore) Geofencing {
	return &geofencing{redis}
}

type locationData struct {
	Ip          string `json:"ip"`
	CountryCode string `json:"country_code"`
	RegionCode  string `json:"region_code"`
	RegionName  string `json:"region_name"`
}

func (g geofencing) IsAllowed(ip string) (bool, error) {
	locationData, err := g.getLocationDataFromRedis(ip)

	/* if not found in cache, get it from the api then set the value to api's response*/
	if err != nil {
		locationData, err = getLocationDataFromAPI(ip)
		if err != nil {
			return false, common.StringError(err)
		}

		err = g.setLocationDataInRedis(ip, locationData)
		if err != nil {
			return false, common.StringError(err)
		}
	}

	return isRegionAllowed(locationData), nil
}

func (c geofencing) setLocationDataInRedis(ip string, locationData locationData) error {
	locationDataStr, err := json.Marshal(locationData)
	if err != nil {
		return common.StringError(err)
	}

	err = c.redis.Set("location-ip"+ip, locationDataStr, A_DAY_IN_NANOSEC)
	if err != nil {
		return common.StringError(err)
	}
	return nil
}

func (g geofencing) getLocationDataFromRedis(ip string) (locationData, error) {
	cachedData, err := g.redis.Get("location-ip" + ip)
	if err != nil {
		return locationData{}, common.StringError(err)
	}

	locationData := locationData{}

	if cachedData == nil {
		return locationData, common.StringError(err)
	}
	err = json.Unmarshal(cachedData, &locationData)
	if err != nil {
		return locationData, common.StringError(err)
	}

	return locationData, nil
}

func getLocationDataFromAPI(ip string) (locationData, error) {
	url := "http://api.ipstack.com/" + ip + "?access_key=" + os.Getenv("LOCATION_API_KEY")

	res, err := http.Get(url)
	if err != nil {
		return locationData{}, common.StringError(err)
	}

	// read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {

		return locationData{}, common.StringError(err)
	}

	dataObj := locationData{}

	// unmarshal the json into our struct
	err = json.Unmarshal(body, &dataObj)
	if err != nil {
		return locationData{}, common.StringError(err)
	}

	if dataObj.Ip != ip || dataObj.CountryCode == "" || dataObj.RegionCode == "" {
		return locationData{}, common.StringError(errors.New(" 4 invalid location data from api"))
	}

	return dataObj, nil
}

func isRegionAllowed(locationData locationData) bool {
	return locationData.CountryCode == "US" && locationData.RegionCode != "NY"
}
