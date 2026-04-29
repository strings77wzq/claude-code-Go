package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type SessionInfo struct {
	ID        string
	FilePath  string
	StartTime time.Time
	EndTime   time.Time
	TurnCount int
	Model     string
	Status    string
}

type listMetaLine struct {
	Type         string `json:"type"`
	SessionID    string `json:"session_id"`
	Model        string `json:"model"`
	StartTime    int64  `json:"start_time_ms"`
	EndTime      int64  `json:"end_time_ms"`
	TurnCount    int    `json:"turn_count"`
	InputTokens  int    `json:"input_tokens"`
	OutputTokens int    `json:"output_tokens"`
	Status       string `json:"status,omitempty"`
}

type listStatusLine struct {
	Type      string `json:"type"`
	Status    string `json:"status"`
	TurnCount int    `json:"turn_count"`
	Timestamp int64  `json:"timestamp_ms"`
}

func ListSessions(dir string) ([]SessionInfo, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("session directory does not exist: %s", dir)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read session directory: %w", err)
	}

	var sessions []SessionInfo

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".jsonl" {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())

		info, err := readSessionMeta(filePath)
		if err != nil {
			continue
		}

		sessions = append(sessions, info)
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].StartTime.After(sessions[j].StartTime)
	})

	return sessions, nil
}

func readSessionMeta(filePath string) (SessionInfo, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return SessionInfo{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	var meta *listMetaLine
	var status *listStatusLine
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		var base struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(line, &base); err != nil {
			continue
		}
		switch base.Type {
		case "meta":
			var parsed listMetaLine
			if err := json.Unmarshal(line, &parsed); err != nil {
				return SessionInfo{}, fmt.Errorf("failed to parse metadata: %w", err)
			}
			meta = &parsed
		case "status":
			var parsed listStatusLine
			if err := json.Unmarshal(line, &parsed); err == nil {
				status = &parsed
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return SessionInfo{}, fmt.Errorf("failed to read file: %w", err)
	}
	if meta == nil {
		return SessionInfo{}, fmt.Errorf("first line is not metadata")
	}

	turnCount := meta.TurnCount
	statusText := meta.Status
	endTime := time.UnixMilli(meta.EndTime)
	if status != nil {
		if status.TurnCount != 0 {
			turnCount = status.TurnCount
		}
		if status.Status != "" {
			statusText = status.Status
		}
		if status.Timestamp != 0 {
			endTime = time.UnixMilli(status.Timestamp)
		}
	}

	return SessionInfo{
		ID:        meta.SessionID,
		FilePath:  filePath,
		StartTime: time.UnixMilli(meta.StartTime),
		EndTime:   endTime,
		TurnCount: turnCount,
		Model:     meta.Model,
		Status:    statusText,
	}, nil
}

// GetLastSessionFilePath returns the file path of the most recent session.
func GetLastSessionFilePath(dir string) (string, error) {
	sessions, err := ListSessions(dir)
	if err != nil {
		return "", err
	}
	if len(sessions) == 0 {
		return "", fmt.Errorf("no sessions found")
	}
	return sessions[0].FilePath, nil
}
