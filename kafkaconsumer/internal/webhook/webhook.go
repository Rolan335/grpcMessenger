package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Caller struct {
	method     string
	url        string
	maxRetries int
	timeout    int
}

func NewCaller(method string, url string, maxRetries int, timeout int) *Caller {
	return &Caller{
		method:     method,
		url:        url,
		maxRetries: maxRetries,
		timeout:    timeout,
	}
}

func (c *Caller) Call(chatUUID string) error {
	for i := 0; i < c.maxRetries; i++ {
		if err := c.call(chatUUID); err != nil {
			if errors.Is(err, ErrWebhookCall) {
				time.Sleep(time.Second)
				fmt.Printf("failed to call webhook for %s at %d retry\n", chatUUID, i+1)
				continue
			}
			return err
		}
		return nil
	}
	return ErrWebhookRetriesExceeded
}

func (c *Caller) call(chatUUID string) error {
	data, err := json.Marshal(map[string]string{"chat_uuid": chatUUID})
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	req, err := http.NewRequest(c.method, c.url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Duration(c.timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ErrWebhookCall
	}

	return nil
}
