package index

import "embed"

//go:embed *.json
var ConfigFiles embed.FS

func LoadJSONFile(filename string) ([]byte, error) {
	return ConfigFiles.ReadFile(filename)
}
