package helper

import "strings"

type AIFormatterResponse struct {
	Prompt_Tokens        uint32
	Completion_Tokens    uint32
	Total_Tokens         uint32
	Processed_Text       string
	Processed_Text_Array []string
}

// ! Cuando tengamos la api key de la ia cambiar la función
func AIFormatter(text_entry string) (*AIFormatterResponse, error) {
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

	// Cuando llamemmos a la ia cambiar esto
	// resp, _ := client.CreateChatCompletion(ctx, req)
	// script := Script{
	// 	ID:               uuid.New(),
	// 	Model:            resp.Model,
	// 	PromptTokens:     uint32(resp.Usage.PromptTokens),
	// 	CompletionTokens: uint32(resp.Usage.CompletionTokens),
	// 	TotalTokens:      uint32(resp.Usage.TotalTokens),
	// }

	aiResponse := AIFormatterResponse{
		Prompt_Tokens:        10,
		Completion_Tokens:    10,
		Total_Tokens:         10,
		Processed_Text:       response,
		Processed_Text_Array: lines,
	}

	return &aiResponse, nil
}
