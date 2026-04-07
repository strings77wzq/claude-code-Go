package telemetry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Client struct {
	enabled bool
	dataDir string
}

type Event struct {
	Timestamp time.Time              `json:"timestamp"`
	EventType string                 `json:"event_type"`
	Data      map[string]interface{} `json:"data"`
}

func NewClient(dataDir string) *Client {
	return &Client{
		enabled: false,
		dataDir: dataDir,
	}
}

func (c *Client) IsEnabled() bool {
	return c.enabled
}

func (c *Client) Enable() error {
	c.enabled = true
	return os.MkdirAll(c.dataDir, 0755)
}

func (c *Client) Disable() {
	c.enabled = false
}

func (c *Client) Track(eventType string, data map[string]interface{}) error {
	if !c.enabled {
		return nil
	}

	event := Event{
		Timestamp: time.Now(),
		EventType: eventType,
		Data:      data,
	}

	filename := filepath.Join(c.dataDir, "telemetry.jsonl")
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(event)
}

func (c *Client) GetEvents() ([]Event, error) {
	filename := filepath.Join(c.dataDir, "telemetry.jsonl")
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []Event{}, nil
		}
		return nil, err
	}

	var events []Event
	lines := splitLines(string(data))
	for _, line := range lines {
		if line == "" {
			continue
		}
		var event Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
