//go:build unit

package validation_test

import (
	"fmt"
	"testing"
	"unicode/utf8"

	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/MetaDandy/cuent-ai-core/src/modules/script"
	"github.com/go-playground/validator/v10"
)

const maxTTSChars = 200

// Helper function para generar strings de longitud específica
func generateString(length int) string {
	result := make([]rune, length)
	for i := range result {
		result[i] = 'a'
	}
	return string(result)
}

func TestLineValidation_TTS_MaxChars(t *testing.T) {
	validate := validator.New()

	// Registrar la validación personalizada igual que en el código real
	validate.RegisterStructValidation(func(sl validator.StructLevel) {
		line := sl.Current().Interface().(script.Line)

		if line.Type == model.AudioTTS &&
			utf8.RuneCountInString(line.Text) > maxTTSChars {
			sl.ReportError(line.Text,
				"Text",
				"text",
				"maxCharsTTS",
				string(rune(maxTTSChars)))
		}
	}, script.Line{})

	tests := []struct {
		name      string
		line      script.Line
		shouldErr bool
	}{
		{
			name: "Valid TTS line - under limit",
			line: script.Line{
				Type: model.AudioTTS,
				Text: "This is a valid TTS line",
			},
			shouldErr: false,
		},
		{
			name: "Valid TTS line - exactly at limit",
			line: script.Line{
				Type: model.AudioTTS,
				Text: generateString(maxTTSChars),
			},
			shouldErr: false,
		},
		{
			name: "Invalid TTS line - over limit by 1",
			line: script.Line{
				Type: model.AudioTTS,
				Text: generateString(maxTTSChars + 1),
			},
			shouldErr: true,
		},
		{
			name: "Invalid TTS line - way over limit",
			line: script.Line{
				Type: model.AudioTTS,
				Text: generateString(maxTTSChars + 100),
			},
			shouldErr: true,
		},
		{
			name: "SFX line - no limit (over TTS limit)",
			line: script.Line{
				Type: model.AudioSFX,
				Text: generateString(maxTTSChars + 100), // Mucho más largo que el límite TTS
			},
			shouldErr: false,
		},
		{
			name: "Empty TTS line - fails required validation",
			line: script.Line{
				Type: model.AudioTTS,
				Text: "",
			},
			shouldErr: true, // Falla porque Text tiene validación "required"
		},
		{
			name: "TTS line - near limit (199 chars)",
			line: script.Line{
				Type: model.AudioTTS,
				Text: generateString(maxTTSChars - 1),
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.line)
			hasErr := err != nil

			if hasErr != tt.shouldErr {
				t.Errorf("validation error mismatch: expected error=%v, got error=%v", tt.shouldErr, hasErr)
				if err != nil {
					t.Logf("validation error details: %v", err)
				}
			}
		})
	}
}

func TestLineValidation_CharacterCount(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{"Empty string", "", 0},
		{"Simple text", "Hello", 5},
		{"With spaces", "Hello World", 11},
		{"Unicode characters", "Hola 世界", 7}, // 4 letras + 2 espacios + 2 caracteres chinos
		{"At limit", generateString(maxTTSChars), maxTTSChars},
		{"Over limit", generateString(maxTTSChars + 1), maxTTSChars + 1},
		{"With newlines", "Line1\nLine2", 11},
		{"With tabs", "Tab\tTab", 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := utf8.RuneCountInString(tt.text)
			if count != tt.expected {
				t.Errorf("character count mismatch: expected %d, got %d", tt.expected, count)
			}
		})
	}
}

func TestLineValidation_TypeComparison(t *testing.T) {
	tests := []struct {
		name     string
		lineType model.AudioLine
		expected bool
	}{
		{"TTS type", model.AudioTTS, true},
		{"SFX type", model.AudioSFX, false},
		{"Empty type", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isTTS := tt.lineType == model.AudioTTS
			if isTTS != tt.expected {
				t.Errorf("type comparison failed: expected %v, got %v", tt.expected, isTTS)
			}
		})
	}
}

func TestLineValidation_MultipleLines(t *testing.T) {
	validate := validator.New()
	validate.RegisterStructValidation(func(sl validator.StructLevel) {
		line := sl.Current().Interface().(script.Line)

		if line.Type == model.AudioTTS &&
			utf8.RuneCountInString(line.Text) > maxTTSChars {
			sl.ReportError(line.Text,
				"Text",
				"text",
				"maxCharsTTS",
				string(rune(maxTTSChars)))
		}
	}, script.Line{})

	lines := []script.Line{
		{Type: model.AudioTTS, Text: "Short line"},
		{Type: model.AudioTTS, Text: generateString(maxTTSChars)}, // At limit
		{Type: model.AudioTTS, Text: generateString(maxTTSChars + 1)}, // Over limit
		{Type: model.AudioSFX, Text: generateString(maxTTSChars + 100)}, // SFX, no limit
	}

	for i, line := range lines {
		t.Run(fmt.Sprintf("Line %d", i+1), func(t *testing.T) {
			err := validate.Struct(line)
			shouldErr := i == 2 // Solo la línea 3 debería fallar
			hasErr := err != nil

			if hasErr != shouldErr {
				t.Errorf("line %d: expected error=%v, got error=%v", i+1, shouldErr, hasErr)
			}
		})
	}
}

func TestLineValidation_BoundaryConditions(t *testing.T) {
	validate := validator.New()
	validate.RegisterStructValidation(func(sl validator.StructLevel) {
		line := sl.Current().Interface().(script.Line)

		if line.Type == model.AudioTTS &&
			utf8.RuneCountInString(line.Text) > maxTTSChars {
			sl.ReportError(line.Text,
				"Text",
				"text",
				"maxCharsTTS",
				string(rune(maxTTSChars)))
		}
	}, script.Line{})

	tests := []struct {
		name      string
		length    int
		shouldErr bool
	}{
		{"Exactly at limit", maxTTSChars, false},
		{"One over limit", maxTTSChars + 1, true},
		{"One under limit", maxTTSChars - 1, false},
		{"Zero characters", 0, true}, // Falla porque Text tiene validación "required"
		{"Very long", maxTTSChars * 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line := script.Line{
				Type: model.AudioTTS,
				Text: generateString(tt.length),
			}

			err := validate.Struct(line)
			hasErr := err != nil

			if hasErr != tt.shouldErr {
				t.Errorf("boundary test failed for length %d: expected error=%v, got error=%v",
					tt.length, tt.shouldErr, hasErr)
			}
		})
	}
}

