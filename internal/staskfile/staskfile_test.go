package staskfile_test

import (
	"path"
	"testing"

	"github.com/itsFrank/stask/internal/staskfile"
	"github.com/stretchr/testify/assert"
)

func TestStaskfileRoundTripHappy(t *testing.T) {
	var tempdir = t.TempDir()

	var tests = []struct {
		name      string
		staskfile staskfile.Staskfile
	}{
		{
			"OneTaskOneState.json",
			staskfile.Staskfile{
				Tasks: map[string]string{"hello": "hello task"},
				State: map[string]string{"state": "foo"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var staskfilePath = path.Join(tempdir, tt.name)
			err := staskfile.WriteStaskfile(staskfilePath, tt.staskfile)
			assert.Nil(t, err)
			stasfkileIn, err := staskfile.ReadStaskfile(staskfilePath)
			assert.Nil(t, err)
			assert.Equal(t, tt.staskfile, stasfkileIn)
		})
	}
}
