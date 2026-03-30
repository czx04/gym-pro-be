package expo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const pushURL = "https://exp.host/--/api/v2/push/send"

// Message is a single Expo push notification.
type Message struct {
	To    string `json:"to"`
	Title string `json:"title,omitempty"`
	Body  string `json:"body"`
	Sound string `json:"sound,omitempty"`
}

type pushRequest struct {
	Messages []Message `json:"messages"`
}

// SendMessages sends up to many messages in one HTTP call (Expo accepts batches).
func SendMessages(ctx context.Context, accessToken string, messages []Message) error {
	if len(messages) == 0 {
		return nil
	}
	if accessToken == "" {
		return fmt.Errorf("expo access token is empty")
	}
	body, err := json.Marshal(pushRequest{Messages: messages})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, pushURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("expo push HTTP %d", resp.StatusCode)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	return nil
}
