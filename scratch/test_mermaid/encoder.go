package mermaid

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// MermaidState represents the JSON schema expected by the Mermaid Live Editor.
type MermaidState struct {
	Code    string         `json:"code"`
	Mermaid map[string]any `json:"mermaid"`
}

// GenerateLiveURL compresses the given Mermaid code and returns a URL to the
// official Mermaid Live Editor pre-loaded with the diagram.
func GenerateLiveURL(code string) (string, error) {
	state := MermaidState{
		Code: code,
		Mermaid: map[string]any{
			"theme": "default",
		},
	}

	jsonBytes, err := json.Marshal(state)
	if err != nil {
		return "", fmt.Errorf("failed to marshal state: %w", err)
	}

	var buf bytes.Buffer
	w, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		return "", fmt.Errorf("failed to initialize zlib: %w", err)
	}

	if _, err := w.Write(jsonBytes); err != nil {
		return "", fmt.Errorf("failed to compress state: %w", err)
	}

	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to finalize compression: %w", err)
	}

	encoded := base64.RawURLEncoding.EncodeToString(buf.Bytes())
	return fmt.Sprintf("https://mermaid.live/edit#pako:%s", encoded), nil
}

// InsertLiveLinks scans markdown content for ```mermaid blocks, generates
// a live URL for each, and appends a clickable link directly beneath the block.
func InsertLiveLinks(content string) string {
	var result strings.Builder
	lines := strings.Split(content, "\n")
	inMermaid := false
	var mermaidCode strings.Builder

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```mermaid") {
			inMermaid = true
			mermaidCode.Reset()
			result.WriteString(line)
			result.WriteByte('\n')
			continue
		}

		if inMermaid {
			if strings.HasPrefix(trimmed, "```") {
				inMermaid = false
				result.WriteString(line)
				result.WriteByte('\n')

				// Generate the URL and append it.
				url, err := GenerateLiveURL(strings.TrimSpace(mermaidCode.String()))
				if err == nil {
					result.WriteString(fmt.Sprintf("\n[Open in Mermaid Live Editor ↗](%s)\n", url))
				}
				continue
			}
			mermaidCode.WriteString(line)
			mermaidCode.WriteByte('\n')
			result.WriteString(line)
			result.WriteByte('\n')
			continue
		}

		result.WriteString(line)
		if i < len(lines)-1 {
			result.WriteByte('\n')
		}
	}

	return result.String()
}
