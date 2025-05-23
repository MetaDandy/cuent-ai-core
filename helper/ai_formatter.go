package helper

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"google.golang.org/genai"
)

type AIFormatterResponse struct {
	Prompt_Tokens        uint32
	Completion_Tokens    uint32
	Total_Tokens         uint32
	Processed_Text       string
	Processed_Text_Array []string
}

func AIFormatter(text_entry string) (*AIFormatterResponse, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	propmt := fmt.Sprintf(`
	You are a text-to-speech pre-processor.
	Task:
	1. Read the user text delimited by triple slash.
	2. Split it into logical lines (one sentence or dialogue unit per line).
	3. Whenever the text mentions or implies a sound effect
	(e.g. shattering glass, footsteps, thunder), output a **separate** line
	that starts with an asterisk (*) followed by a detailed English,
	onomatopoeic description of the sound suitable for ElevenLabs
	(Example: *shattering glass — sharp crystalline crack followed by tinkling fragments scattering on a hard tile floor*).
	4. Keep the original language for normal narrative or dialogue lines;
	only the sound-effect lines must be in English.
	5. Return **only** the processed lines, one per output line, with no extra commentary.
	///%s///
	`, text_entry)

	resp, err := client.Models.GenerateContent(
		ctx,
		os.Getenv("GEMINI_MODEL"),
		genai.Text(propmt),
		nil,
	)
	if err != nil {
		return nil, err
	}

	out := resp.Text()

	var lines []string
	for line := range strings.SplitSeq(out, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}

	usage := resp.UsageMetadata
	aiResponse := AIFormatterResponse{
		Prompt_Tokens:        uint32(usage.PromptTokenCount),
		Completion_Tokens:    uint32(usage.CandidatesTokenCount),
		Total_Tokens:         uint32(usage.TotalTokenCount),
		Processed_Text:       out,
		Processed_Text_Array: lines,
	}

	return &aiResponse, nil
}

// ! Función de emergencia buscar una solución mejor.
// EstimateCuentokens devuelve una cota superior del número de líneas
// (≈ cuentokens) que obtendrías tras formatear `text`.
//
// Estrategia:
//  1. Cuenta las frases terminadas en ".", "!", "?".
//  2. Añade las líneas ya existentes (separadas por \n) que NO terminan
//     en esos signos (p. ej. diálogos separados con guion largo —).
//  3. Añade un 10 % extra como colchón para posibles líneas SFX.
func EstimateCuentokens(text string) uint {
	// Paso 1: frases que terminan en . ! ?
	re := regexp.MustCompile(`[^.!?]+[.!?]`)
	frases := re.FindAllString(text, -1)
	nFrases := len(frases)

	// Paso 2: líneas “sueltas” (diálogos / listas) que no detectó el regex
	nuevasLineas := 0
	for _, l := range strings.Split(text, "\n") {
		trim := strings.TrimSpace(l)
		if trim == "" {
			continue
		}
		if !re.MatchString(trim) { // no termina en . ! ?
			nuevasLineas++
		}
	}

	total := nFrases + nuevasLineas

	// Paso 3: 10 % extra para sonidos (*SFX*)
	extra := (total + 9) / 10 // redondeo hacia arriba
	return uint(total + extra)
}
