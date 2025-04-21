package cuentai

import (
	"fmt"
	"strings"

	tts "github.com/MetaDandy/cuent-ai-core/src/modules/Tts"
)

type Service struct {
	tts *tts.Service
}

func NewService(tts *tts.Service) *Service {
	return &Service{
		tts: tts,
	}
}

// CuentAIFlow formatea, genera TTS/SFX y concatena audio
func (s *Service) CuentAIFlow(textEntry string) (*CuentAIFlowResult, error) {
	// 1) formatear
	lines, err := s.aiFormatter(textEntry)
	if err != nil {
		return nil, fmt.Errorf("AIFormatter error: %w", err)
	}
	// 2) generar audio
	ttsClips, sfxClips, combined, err := s.audioOutput(lines)
	if err != nil {
		return nil, fmt.Errorf("AudioOutput error: %w", err)
	}
	// 3) devolver todo
	return &CuentAIFlowResult{
		Lines:    lines,
		TTSClips: ttsClips,
		SFXClips: sfxClips,
		Combined: combined,
	}, nil
}

func (s *Service) aiFormatter(text_entry string) ([]string, error) {
	// Aquí simulas la respuesta de la API tal como en el ejemplo de Python
	response := `SACERDOTE.- Con oportunidad has hablado
	Precisamente éstos me están indicando por señas que Creonte se acerca
	* Entrance of Creon
	EDIPO.- ¡Oh soberano Apolo! ¡Ojalá viniera con suerte liberadora, del mismo modo que viene con rostro radiante!
	SACERDOTE.- Por lo que se puede adivinar, viene complacido
	En otro caso no vendría así, con la cabeza coronada de frondosas ramas de laurel
	EDIPO.- Pronto lo sabremos, pues ya está lo suficientemente cerca para que nos escuche
	¡Oh príncipe, mi pariente, hijo de Meneceo!
	¿Cuál es la respuesta del oráculo?
	CREONTE.- Con una buena
	Afirmo que incluso las aflicciones, si llegan felizmente a término, todas pueden resultar bien
	EDIPO.- ¿Cuál es la respuesta?
	Por lo que acabas de decir, no estoy ni tranquilo ni tampoco preocupado
	CREONTE.- Si deseas oírlo estando éstos aquí cerca, estoy dispuesto a hablar y también, si lo deseas, a ir dentro
	EDIPO.- Habla ante todos, ya que por ellos sufro una aflicción mayor, incluso, que por mi propia vida
	CREONTE.- Diré las palabras que escuché de parte del dios
	El soberano Febo nos ordenó, claramente, arrojar de la región una mancilla que existe en esta tierra y no mantenerla para que llegue a ser irremediable
	EDIPO.- ¿Con qué expiación?
	¿Cuál es la naturaleza de la desgracia?
	CREONTE.- Con el destierro o liberando un antiguo asesinato con otro, puesto que esta sangre es la que está sacudiendo la ciudad
	EDIPO.- ¿De qué hombre denuncia tal desdicha?
	CREONTE.- Teníamos nosotros, señor, en otro tiempo a Layo como soberano de esta tierra, antes de que tú rigieras rectamente esta ciudad
	EDIPO.- Lo sé por haberlo oído, pero nunca lo vi`

	var lines []string
	for _, line := range strings.Split(response, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}

	return lines, nil
}

func (s *Service) audioOutput(
	lines []string,
) (ttsClips [][]byte, sfxClips [][]byte, global []byte, err error) {
	for idx, raw := range lines {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "*") {
			// Esto es un SFX: quitamos el '*' y el posible espacio
			prompt := strings.TrimSpace(strings.TrimPrefix(line, "*"))
			audio, err := s.tts.TextToSoundEffects(
				prompt,
				3.0,
				1.0,
				"mp3_44100_128",
			)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("SFX falló en línea %d: %w", idx, err)
			}
			sfxClips = append(sfxClips, audio)
			global = append(global, audio...)
		} else {
			// Esto es TTS normal
			audio, err := s.tts.TextToSpeechElevenlabs(line, "")
			if err != nil {
				return nil, nil, nil, fmt.Errorf("TTS falló en línea %d: %w", idx, err)
			}
			ttsClips = append(ttsClips, audio)
			global = append(global, audio...)
		}
	}
	return ttsClips, sfxClips, global, nil
}
