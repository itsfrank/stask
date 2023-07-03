package staskfile

import (
	"encoding/json"
	"os"
)

type Staskfile struct {
	Tasks map[string]string
	State map[string]string
}

func ParseStaskfile(data []byte) (Staskfile, error) {
	var staskfile Staskfile
	err := json.Unmarshal(data, &staskfile)
	if err != nil {
		return Staskfile{}, nil
	}
	return staskfile, nil
}

func SerializeStaskfile(staskfile Staskfile) ([]byte, error) {
	return json.MarshalIndent(staskfile, "", "    ")
}

func ReadStaskfile(path string) (Staskfile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Staskfile{}, err
	}
	return ParseStaskfile(data)
}
func WriteStaskfile(path string, staskfile Staskfile) error {
	data, err := SerializeStaskfile(staskfile)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0777)
}
