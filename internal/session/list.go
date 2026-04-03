package session

import (
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

	var firstLine string
	buf := make([]byte, 1)
	for {
		n, err := f.Read(buf)
		if err != nil {
			return SessionInfo{}, fmt.Errorf("failed to read file: %w", err)
		}
		if n == 0 {
			break
		}
		if buf[0] == '\n' {
			break
		}
		firstLine += string(buf[0])
	}

	if firstLine == "" {
		return SessionInfo{}, fmt.Errorf("empty file")
	}

	var meta listMetaLine
	if err := json.Unmarshal([]byte(firstLine), &meta); err != nil {
		return SessionInfo{}, fmt.Errorf("failed to parse metadata: %w", err)
	}

	if meta.Type != "meta" {
		return SessionInfo{}, fmt.Errorf("first line is not metadata")
	}

	return SessionInfo{
		ID:        meta.SessionID,
		FilePath:  filePath,
		StartTime: time.UnixMilli(meta.StartTime),
		EndTime:   time.UnixMilli(meta.EndTime),
		TurnCount: meta.TurnCount,
		Model:     meta.Model,
	}, nil
}
