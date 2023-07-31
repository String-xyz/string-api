package checkout

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/String-xyz/string-api/config"
)

type SlackMessage struct {
	Text string `json:"text"`
}

// This function sends a message to a Slack channel via Incoming Webhooks
func PostToSlack(event EventType) {
	webhookURL := config.Var.SLACK_WEBHOOK_URL
	msg := SlackMessage{
		Text: "Event of type " + string(event) + " received.",
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Error().Msgf("Error marshalling Slack message: %s", err.Error())
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(msgBytes))
	if err != nil {
		log.Error().Msgf("Error posting to Slack: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Error().Msgf("Received non-2xx response code: %d", resp.StatusCode)
	}
}
