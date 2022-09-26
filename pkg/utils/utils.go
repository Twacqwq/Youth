package utils

import (
	"encoding/json"
	"os"
	"path/filepath"

	pkg "github.com/Twacqwq/youth/pkg/youth"
)

func LoadConfig(configDir string) (*os.File, error) {
	configPath := filepath.Join(configDir, "config.json")
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func GenerateConfig() {
	filename := "config.json"
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()
	exampleConfig := []pkg.Member{
		{MemberId: 123},
		{MemberId: 123},
		{MemberId: 123},
	}
	enc := json.NewEncoder(f)
	enc.Encode(exampleConfig)
}
