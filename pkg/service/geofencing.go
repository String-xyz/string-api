package service

import (
	"encoding/json"
	"io"
	"net/http"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	"github.com/String-xyz/string-api/env"
	"github.com/pkg/errors"
)

const A_DAY_IN_NANOSEC = 84600000000000

type Geofencing interface {
	IsAllowed(ip string) (bool, error)
	getLocation(ip string) (GeoLocation, error)
	setLocation(ip string, location GeoLocation) error
}

type geofencing struct {
	redis database.RedisStore
}

func NewGeofencing(redis database.RedisStore) Geofencing {
	return &geofencing{redis}
}

type GeoLocation struct {
	Ip          string `json:"ip"`
	CountryCode string `json:"country_code"`
	RegionCode  string `json:"region_code"`
	RegionName  string `json:"region_name"`
}

func (g geofencing) IsAllowed(ip string) (bool, error) {
	// location, err := g.getLocation(ip)

	// /* if not found in cache, get it from the api then set the value to api's response*/
	// if err != nil {
	// 	location, err = getLocationFromAPI(ip)
	// 	if err != nil {
	// 		return false, libcommon.StringError(err)
	// 	}

	// 	err = g.setLocation(ip, location)
	// 	if err != nil {
	// 		return false, libcommon.StringError(err)
	// 	}
	// }

	// return isRegionAllowed(location), nil
	return true, nil
}

func (c geofencing) setLocation(ip string, location GeoLocation) error {
	locationStr, err := json.Marshal(location)
	if err != nil {
		return libcommon.StringError(err)
	}

	err = c.redis.Set("location-ip"+ip, locationStr, A_DAY_IN_NANOSEC)
	if err != nil {
		return libcommon.StringError(err)
	}
	return nil
}

func (g geofencing) getLocation(ip string) (GeoLocation, error) {
	cachedData, err := g.redis.Get("location-ip" + ip)
	if err != nil {
		return GeoLocation{}, libcommon.StringError(err)
	}

	location := GeoLocation{}

	if cachedData == nil {
		return location, libcommon.StringError(err)
	}
	err = json.Unmarshal(cachedData, &location)
	if err != nil {
		return location, libcommon.StringError(err)
	}

	return location, nil
}

func getLocationFromAPI(ip string) (GeoLocation, error) {
	key, err := env.Get("IPSTACK_API_KEY")
	if err != nil {
		return GeoLocation{}, libcommon.StringError(err)
	}
	url := "http://api.ipstack.com/" + ip + "?access_key=" + key

	res, err := http.Get(url)
	if err != nil {
		return GeoLocation{}, libcommon.StringError(err)
	}

	// read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {

		return GeoLocation{}, libcommon.StringError(err)
	}

	dataObj := GeoLocation{}

	// unmarshal the json into our struct
	err = json.Unmarshal(body, &dataObj)
	if err != nil {
		return GeoLocation{}, libcommon.StringError(err)
	}

	if dataObj.Ip != ip || dataObj.CountryCode == "" || dataObj.RegionCode == "" {
		return GeoLocation{}, libcommon.StringError(errors.New("The Data returned by the external location service is invalid"))
	}

	return dataObj, nil
}

func isRegionAllowed(location GeoLocation) bool {
	return location.CountryCode == "US" && location.RegionCode != "NY"
}
