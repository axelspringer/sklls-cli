package tracking

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sklls/api"
)

type TrackingFn func(sessionId string, eventType string, userId string, metaJson string) error

func GetEventTrackingFn(apiUrl string, apiKey string) TrackingFn {
	trackingUrl := apiUrl + "/tracking"

	return func(sessionId string, eventType string, userId string, metaJson string) error {
		addTrackingEventRequest := api.TrackingEventCreationRequest{
			SessionId: sessionId,
			ApiKey:    apiKey,
			Type:      eventType,
			UserId:    userId,
			MetaJson:  metaJson,
		}
		requestBodyJson, err := json.Marshal(addTrackingEventRequest)

		// Create profile & parse response
		req, err := http.NewRequest("POST", trackingUrl, bytes.NewBuffer(requestBodyJson))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		rawJsonBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		parsedRes := api.ProfileCreationResponse{}
		err = json.Unmarshal([]byte(rawJsonBody), &parsedRes)
		if err != nil {
			return err
		}

		log.Printf("[TRACKING] %s - %s - %s\n", userId, eventType, metaJson)
		return nil
	}
}
