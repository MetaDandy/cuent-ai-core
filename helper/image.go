package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// Image contiene los datos que quieres extraer de cada resultado.
// Aquí guardamos sólo la URL de tamaño "regular", pero puedes añadir más campos.
type Image struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}

// searchResponse refleja la estructura parcial de la respuesta JSON de Unsplash.
type searchResponse struct {
	Results []struct {
		ID   string `json:"id"`
		Urls struct {
			Regular string `json:"regular"`
		} `json:"urls"`
	} `json:"results"`
}

// SearchImage consulta la API de Unsplash para obtener todas las imágenes
// relacionadas con el texto `prompt`. Devuelve un slice de Image o error.
func SearchImage(prompt string) ([]Image, error) {
	apiKey := os.Getenv("UNSPLASH_ACCESS_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("debes definir la variable de entorno UNSPLASH_ACCESS_KEY")
	}

	endpoint := "https://api.unsplash.com/search/photos"
	params := url.Values{}
	params.Set("client_id", apiKey)
	params.Set("query", prompt)
	params.Set("page", "1")
	params.Set("per_page", "10")

	req, err := http.NewRequest("GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error Unsplash API: code %d", resp.StatusCode)
	}

	var data searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	images := make([]Image, len(data.Results))
	for i, r := range data.Results {
		images[i] = Image{
			ID:  r.ID,
			Url: r.Urls.Regular,
		}
	}
	return images, nil
}
