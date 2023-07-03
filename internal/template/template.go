package template

import (
	"errors"
	"strings"
)

type Key struct {
	Index int
	Str   string
}

type Template struct {
	Str  string
	Keys []Key
}

func ParseTemplate(str string) (Template, error) {
	template := Template{}
	curKey := ""
	curKeyIndex := 0

	index := 0
	activeKey := false
	var builder strings.Builder
	builder.Grow(len(str))
	for _, chr := range str {
		switch chr {

		case '{':
			if activeKey {
				return Template{}, errors.New("Found opening '{' before closing '}'")
			}
			curKeyIndex = index
			activeKey = true
			continue
		case '}':
			if activeKey && len(curKey) == 0 {
				return Template{}, errors.New("Empty Key")
			}
			if !activeKey {
				return Template{}, errors.New("Found closing '}' before opening '{'")
			}
			template.Keys = append(template.Keys, Key{curKeyIndex, curKey})
			activeKey = false
			curKey = ""
			continue
		default:
			if activeKey {
				curKey += string(chr)
				continue
			}
		}
		index += 1
		builder.WriteRune(chr)
	}
	template.Str = builder.String()
	return template, nil
}

// returns the string with values applied, and a list of keys that were not found in the input map
// most applications would consider len(missing) > 0 to be an error
func ApplyTemplate(tmpl Template, values map[string]string) (string, []string) {
	var missing []string
	var lastIndex int
	var builder strings.Builder
	for _, key := range tmpl.Keys {
		val, prs := values[key.Str]
		if !prs {
			missing = append(missing, key.Str)
			continue
		}

		builder.WriteString(tmpl.Str[lastIndex:key.Index])
		builder.WriteString(val)
		lastIndex = key.Index
	}
	builder.WriteString(tmpl.Str[lastIndex:])
	return builder.String(), missing
}
