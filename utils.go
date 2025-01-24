package harlyzer

import (
	"encoding/json"
	"fmt"
	"os"
)

func ParseHarFile(path string) (*HAR, error) {
	var har HAR
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open HAR file: %w", err)
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&har)
	if err != nil {
		return nil, fmt.Errorf("failed to decode HAR file: %w", err)
	}
	return &har, nil
}
