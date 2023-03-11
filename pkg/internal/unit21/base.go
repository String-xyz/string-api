package unit21

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/rs/zerolog/log"
)

func u21Put(url string, jsonBody any) (body []byte, err error) {
	apiKey := os.Getenv("UNIT21_API_KEY")

	reqBodyBytes, err := json.Marshal(jsonBody)
	if err != nil {
		log.Err(err).Msg("Could not encode into bytes")
		return nil, common.StringError(err)
	}
	log.Info().Str("body", string(reqBodyBytes)).Send()
	bodyReader := bytes.NewReader(reqBodyBytes)

	req, err := http.NewRequest(http.MethodPut, url, bodyReader)
	if err != nil {
		log.Err(err).Str("url", url).Msg("Could not create request")
		return nil, common.StringError(err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("u21-key", apiKey)

	client := http.Client{Timeout: 10 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		log.Err(err).Str("url", url).Msg("Request failed to update")
		return nil, common.StringError(err)
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		log.Err(err).Str("url", url).Msg("Error extracting body")
		return nil, common.StringError(err)
	}

	if res.StatusCode != 200 {
		log.Err(err).Str("url", url).Int("statusCode", res.StatusCode).Msg("Request failed to update")
		err = common.StringError(fmt.Errorf("request failed with status code %s and return body: %s", fmt.Sprint(res.StatusCode), string(body)))
		return
	}

	return body, nil
}

func u21Post(url string, jsonBody any) (body []byte, err error) {
	apiKey := os.Getenv("UNIT21_API_KEY")

	reqBodyBytes, err := json.Marshal(jsonBody)
	if err != nil {
		log.Err(err).Msg("Could not encode into bytes")
		return nil, common.StringError(err)
	}

	bodyReader := bytes.NewReader(reqBodyBytes)

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		log.Err(err).Str("url", url).Msg("Could not create request")
		return nil, common.StringError(err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("u21-key", apiKey)

	client := http.Client{Timeout: 10 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		log.Err(err).Str("url", url).Msg("Request failed to update")
		return nil, common.StringError(err)
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		log.Err(err).Str("url", url).Msg("Error extracting body from")
		return nil, common.StringError(err)
	}

	log.Info().Str("body", string(body)).Msgf("String of body from response")

	if res.StatusCode != 200 {
		log.Err(err).Str("url", url).Int("statusCode", res.StatusCode).Msg("Request failed to update")
		err = common.StringError(fmt.Errorf("request failed with status code %s and return body: %s", fmt.Sprint(res.StatusCode), string(body)))
		return
	}

	return body, nil
}
