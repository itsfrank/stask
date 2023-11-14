// adapted from: https://adrianhesketh.com/2017/03/18/getting-line-and-character-positions-from-gos-json-unmarshal-errors
package jsonerror

import (
	"encoding/json"
	"errors"
	"fmt"
)

func GetFormattedError(jsonString string, err error) error {

	if jsonError, ok := err.(*json.SyntaxError); ok {
		line, character, lcErr := lineAndCharacter(jsonString, int(jsonError.Offset))

		if lcErr != nil {
			return err
		}
		return errors.New(fmt.Sprintf("syntax error - line %d, character %d: %v\n", line, character, jsonError.Error()))
	}
	if jsonError, ok := err.(*json.UnmarshalTypeError); ok {
		line, character, lcErr := lineAndCharacter(jsonString, int(jsonError.Offset))
		if lcErr != nil {
			return err
		}
		return errors.New(fmt.Sprintf("json type '%v' cannot be converted into go '%v' type - line %d, character %d\n", jsonError.Value, jsonError.Type.Name(), line, character))
	}

	return err
}

func lineAndCharacter(input string, offset int) (line int, character int, err error) {
	lf := rune(0x0A)

	if offset > len(input) || offset < 0 {
		return 0, 0, fmt.Errorf("Couldn't find offset %d within the input.", offset)
	}

	// Humans tend to count from 1.
	line = 1

	for i, b := range input {
		if b == lf {
			line++
			character = 0
		}
		character++
		if i == offset {
			break
		}
	}

	return line, character, nil
}
