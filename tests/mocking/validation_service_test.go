//go:build mocking
// +build mocking

package mocking

import (
	"strings"
	"testing"

	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/MetaDandy/cuent-ai-core/src/modules/script"
	"github.com/MetaDandy/cuent-ai-core/src/modules/validation"
	"github.com/stretchr/testify/assert"
)

func TestValidation_LineStruct(t *testing.T) {
	tests := []struct {
		name      string
		line      script.Line
		expectErr bool
	}{
		{
			name: "valid TTS within limit",
			line: script.Line{Type: model.AudioTTS, Text: "hola"},
		},
		{
			name:      "too long TTS",
			line:      script.Line{Type: model.AudioTTS, Text: strings.Repeat("a", 205)},
			expectErr: true,
		},
		{
			name: "SFX ignores length",
			line: script.Line{Type: model.AudioSFX, Text: strings.Repeat("b", 500)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.Validate.Struct(tt.line)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
