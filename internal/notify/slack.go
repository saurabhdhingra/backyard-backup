package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SlackPayload struct {
	Text string `json:"text"`
}

func SendSlackNotification(webhookURL, message string) error {
	fmt.Printf("DEBUG: Sending to %s...\n", webhookURL)
	payload := SlackPayload{Text: message}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	fmt.Printf("DEBUG: Payload: %s\n", string(data))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(webhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	fmt.Printf("DEBUG: Slack Response Code: %d\n", resp.StatusCode)
	fmt.Printf("DEBUG: Slack Response Body: %s\n", buf.String())

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack notification failed: %s", buf.String())
	}

	return nil
}
