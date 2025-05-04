package cuentai

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tts "github.com/MetaDandy/cuent-ai-core/src/modules/Tts"
	"github.com/MetaDandy/cuent-ai-core/src/modules/supabase"
)

type Service struct {
	tts     *tts.Service
	storage *supabase.Service
}

func NewService(tts *tts.Service, storage *supabase.Service) *Service {
	return &Service{
		tts:     tts,
		storage: storage,
	}
}

/* TODO:
- Cambiar el nombre del mixtape de forma generica
- Refactorizar las funciones
- Implementar guardado en la db
- Implementar funciones auxiliares, dtos
*/

// CuentAIFlow formatea, genera TTS/SFX y concatena audio
func (s *Service) CuentAIFlow(textEntry string) (*CuentAIFlowResult, error) {
	// 1) formatear
	lines, err := s.aiFormatter(textEntry)
	if err != nil {
		return nil, fmt.Errorf("AIFormatter error: %w", err)
	}
	// 2) generar audio
	ttsURLs, sfxURLs, mixedURL, err := s.audioOutput(lines)
	if err != nil {
		return nil, fmt.Errorf("AudioOutput error: %w", err)
	}
	// 3) devolver todo
	return &CuentAIFlowResult{
		Lines:    lines,
		TTSURLS:  ttsURLs,
		SFXURLS:  sfxURLs,
		MixedURL: mixedURL,
	}, nil
}

func (s *Service) aiFormatter(text_entry string) ([]string, error) {
	// Aquí simulas la respuesta de la API tal como en el ejemplo de Python
	response := `SACERDOTE.- Con oportunidad has hablado
	Precisamente éstos me están indicando por señas que Creonte se acerca
	* Entrance of Creon
	EDIPO.- ¡Oh soberano Apolo! ¡Ojalá viniera con suerte liberadora, del mismo modo que viene con rostro radiante!
	`

	/**

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
	EDIPO.- Lo sé por haberlo oído, pero nunca lo vi
	*/

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
) (ttsURLs []string, sfxURLs []string, mixedURL string, err error) {

	const (
		bucket  string = "audio"
		dirPath string = "user/path"
	)

	var (
		tmpFiles []string    // rutas de los .mp3 temporales
		cleanup  = func() {} // se redefine al final para borrar temporales
	)
	// ↳ defiere la limpieza al final de la función
	defer func() { cleanup() }()

	for idx, raw := range lines {
		line := strings.TrimSpace(raw)
		var audio []byte

		if strings.HasPrefix(line, "*") {
			prompt := strings.TrimSpace(strings.TrimPrefix(line, "*"))
			audio, err = s.tts.TextToSoundEffects(
				prompt,
				3.0,
				1.0,
				"mp3_44100_128",
			)
			if err != nil {
				return nil, nil, "", fmt.Errorf("SFX falló en línea %d: %w", idx, err)
			}

			fileName := fmt.Sprintf("sfx_%d.mp3", idx)
			if url, err := s.storage.Upload(
				context.TODO(), // o el ctx que uses en el flujo completo
				bucket,
				dirPath,
				fileName,
				bytes.NewReader(audio),
				"audio/mpeg",
				true,
			); err != nil {
				return nil, nil, "", fmt.Errorf("Upload SFX %d: %w", idx, err)
			} else {
				sfxURLs = append(sfxURLs, url) // guarda URL
			}
		} else {
			// Esto es TTS normal
			audio, err = s.tts.TextToSpeechElevenlabs(line, "")
			if err != nil {
				return nil, nil, "", fmt.Errorf("TTS falló en línea %d: %w", idx, err)
			}

			fileName := fmt.Sprintf("tts_%d.mp3", idx)
			if url, err := s.storage.Upload(
				context.TODO(),
				bucket,
				dirPath,
				fileName,
				bytes.NewReader(audio),
				"audio/mpeg",
				true,
			); err != nil {
				return nil, nil, "", fmt.Errorf("Upload TTS %d: %w", idx, err)
			} else {
				ttsURLs = append(ttsURLs, url)
			}
		}

		tmpf, errTmp := os.CreateTemp("", fmt.Sprintf("clip_%d_*.mp3", idx))
		if errTmp != nil {
			err = fmt.Errorf("tmpfile: %w", errTmp)
			return
		}
		if _, errTmp = tmpf.Write(audio); errTmp != nil {
			tmpf.Close()
			err = fmt.Errorf("write tmp: %w", errTmp)
			return
		}
		tmpf.Close()
		tmpFiles = append(tmpFiles, tmpf.Name())
	}

	if len(tmpFiles) > 0 {
		// Crear lista para el demuxer concat
		listFile, errList := os.CreateTemp("", "ffmpeg_list_*.txt")
		if errList != nil {
			err = fmt.Errorf("tmp list: %w", errList)
			return
		}
		for _, p := range tmpFiles {
			// Cada línea: file '/path/to/file'
			fmt.Fprintf(listFile, "file '%s'\n", filepath.ToSlash(p))
		}
		listFile.Close()

		// Archivo destino temporal
		mixPath := filepath.Join(os.TempDir(),
			fmt.Sprintf("mix_%d.mp3", time.Now().UnixNano()))

		cmd := exec.Command(
			"ffmpeg", "-y",
			"-f", "concat", "-safe", "0",
			"-i", listFile.Name(),
			"-acodec", "libmp3lame",
			"-b:a", "192k", // 192 kb/s ≈ transparente p/voz y SFX
			"-ar", "44100", // 44,1 kHz uniforme
			"-ac", "2", // estéreo
			mixPath,
		)
		if out, errFfmpeg := cmd.CombinedOutput(); errFfmpeg != nil {
			err = fmt.Errorf("ffmpeg: %v – %s", errFfmpeg, string(out))
			return
		}

		// Leer mix.mp3 y subir
		mixBytes, errRead := os.ReadFile(mixPath)
		if errRead != nil {
			err = fmt.Errorf("read mix: %w", errRead)
			return
		}

		mixName := "mix_" + time.Now().Format("20060102_150405") + ".mp3"
		if mixedURL, err = s.storage.Upload(
			context.TODO(), bucket, dirPath, mixName,
			bytes.NewReader(mixBytes), "audio/mpeg", true,
		); err != nil {
			err = fmt.Errorf("upload mix: %w", err)
			return
		}

		// Redefine cleanup para borrar todos los temporales (clips + list + mix)
		cleanup = func() {
			os.Remove(mixPath)
			os.Remove(listFile.Name())
			for _, f := range tmpFiles {
				_ = os.Remove(f)
			}
		}
	}

	return ttsURLs, sfxURLs, mixedURL, nil
}
