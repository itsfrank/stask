package template_test

import (
	"errors"
	"testing"

	"github.com/itsFrank/stask/internal/template"
	"github.com/stretchr/testify/assert"
)

func TestParseTemplateHappy(t *testing.T) {
	var tests = []struct {
		str      string
		stripped string
		keys     []template.Key
	}{
		{
			"my {adjective} template!",
			"my  template!",
			[]template.Key{{3, "adjective"}},
		},
		{
			"more {than} one {tmpl}",
			"more  one ",
			[]template.Key{{5, "than"}, {10, "tmpl"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			tmpl, err := template.ParseTemplate(tt.str)
			assert.Nil(t, err)
			assert.Equal(t, tt.keys, tmpl.Keys)
			assert.Equal(t, tt.stripped, tmpl.Str)
		})
	}
}

func TestParseTemplateError(t *testing.T) {
	var tests = []struct {
		str string
		err error
	}{
		{
			"my {a{djective} template!",
			errors.New("Found opening '{' before closing '}'"),
		},
		{
			"my {} template!",
			errors.New("Empty Key"),
		},
		{
			"my } template!",
			errors.New("Found closing '}' before opening '{'"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			_, err := template.ParseTemplate(tt.str)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestApplyTemplateHappy(t *testing.T) {
	var tests = []struct {
		str    string
		values map[string]string
		result string
	}{
		{
			"my {adjective} template!",
			map[string]string{"adjective": "cool"},
			"my cool template!",
		},
		{
			"Hi {you} I am {me}, nice to meet you",
			map[string]string{"you": "Mark", "me": "Frank"},
			"Hi Mark I am Frank, nice to meet you",
		},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			tmpl, _ := template.ParseTemplate(tt.str)
			res, mss := template.ApplyTemplate(tmpl, tt.values)
			assert.Equal(t, tt.result, res)
			assert.Equal(t, []string(nil), mss)
		})
	}
}

func TestApplyTemplateError(t *testing.T) {
	var tests = []struct {
		str     string
		values  map[string]string
		missing []string
	}{
		{
			"my {adjective} template!",
			map[string]string{},
			[]string{"adjective"},
		},
		{
			"Hi {you} I am {me}, nice to meet you",
			map[string]string{"them": "Mark", "us": "Frank"},
			[]string{"you", "me"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			tmpl, _ := template.ParseTemplate(tt.str)
			_, mss := template.ApplyTemplate(tmpl, tt.values)
			assert.Equal(t, tt.missing, mss)
		})
	}
}
