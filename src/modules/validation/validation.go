package validation

import (
	"strconv"
	"unicode/utf8"

	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/MetaDandy/cuent-ai-core/src/modules/script"
	"github.com/go-playground/validator/v10"
)

const maxTTSChars = 200

var Validate *validator.Validate

func init() {
	Validate = validator.New()
	Validate.RegisterStructValidation(lineStructValidation, script.Line{})
}

func lineStructValidation(sl validator.StructLevel) {
	line := sl.Current().Interface().(script.Line)

	// Sólo aplica el límite si la línea es TTS
	if line.Type == model.AudioTTS &&
		utf8.RuneCountInString(line.Text) > maxTTSChars {

		sl.ReportError(line.Text,
			"Text",                    // nombre Go del campo
			"text",                    // nombre JSON
			"maxCharsTTS",             // etiqueta de error
			strconv.Itoa(maxTTSChars)) // parámetro extra
	}
}
