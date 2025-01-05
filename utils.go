package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// CreateEntityRef creates an entity reference from an area and ID.
func CreateEntityRef(area, id string) string {
	return strings.ToLower(fmt.Sprintf("%s:%s", area, id))
}

// ParseEntityRef parses an entity reference into its area and ID parts.
func ParseEntityRef(entityRef string) (area, id string) {
	parts := strings.Split(strings.ToLower(entityRef), ":")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// WrapText splits text into lines of the specified width without breaking words.
func WrapText(text string, width int) string {
	var result strings.Builder
	words := strings.Fields(text)
	line := ""

	for _, word := range words {
		if len(line)+len(word)+1 > width {
			result.WriteString(line + "\n")
			line = word
		} else {
			if line != "" {
				line += " "
			}
			line += word
		}
	}
	if line != "" {
		result.WriteString(line)
	}

	return result.String()
}

func LoadYAML(filePath string, out interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	return decoder.Decode(out)
}

func SaveYAML(filePath string, in interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	return encoder.Encode(in)
}

func LoadJSON(filePath string, out interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		slog.Error("failed reading file",
			slog.String("file", filePath),
			slog.Any("error", err))

		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(out)
}

func SaveJSON(filePath string, in interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		slog.Error("failed reading file",
			slog.String("file", filePath),
			slog.Any("error", err))

		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(in)
}

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}
