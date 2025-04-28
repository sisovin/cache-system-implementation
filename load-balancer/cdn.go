package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type PurgeRequest struct {
	URLs []string `json:"urls"`
}

func purgeCDNCache(urls []string) error {
	purgeRequest := PurgeRequest{URLs: urls}
	requestBody, err := json.Marshal(purgeRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal purge request: %v", err)
	}

	req, err := http.NewRequest("POST", "https://cdn.example.com/purge", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create purge request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send purge request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("purge request failed with status: %v", resp.Status)
	}

	return nil
}
