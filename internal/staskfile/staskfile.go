package staskfile

import (
	"encoding/json"
	"github.com/itsFrank/stask/internal/jsonerror"
	"os"
)

type Staskfile struct {
	Tasks    map[string]string
	State    map[string]string
	Profiles map[string]map[string]string
}

func Empty() Staskfile {
	sf := Staskfile{}
	sf.Tasks = map[string]string{}
	sf.State = map[string]string{}
	sf.Profiles = map[string]map[string]string{}
	return sf
}

func ParseStaskfile(data []byte) (Staskfile, error) {
	var staskfile Staskfile
	err := json.Unmarshal(data, &staskfile)

	if err != nil {
		return Empty(), jsonerror.GetFormattedError(string(data), err)
	}

	// make sure all maps are initialized
	if staskfile.Tasks == nil {
		staskfile.Tasks = map[string]string{}
	}
	if staskfile.State == nil {
		staskfile.State = map[string]string{}
	}
	if staskfile.Profiles == nil {
		staskfile.Profiles = map[string]map[string]string{}
	}

	return staskfile, nil
}

func SerializeStaskfile(staskfile Staskfile) ([]byte, error) {
	return json.MarshalIndent(staskfile, "", "    ")
}

func ReadStaskfile(path string) (Staskfile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Empty(), err
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
