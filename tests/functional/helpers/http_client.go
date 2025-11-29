//go:build e2e

package helpers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "testing"
    "time"
)

// HTTPClient es un cliente HTTP reutilizable para pruebas E2E
type HTTPClient struct {
    BaseURL    string
    HTTPClient *http.Client
    AuthToken  string
    T          *testing.T
}

// NewHTTPClient crea un nuevo cliente HTTP para pruebas
func NewHTTPClient(baseURL string, t *testing.T) *HTTPClient {
    return &HTTPClient{
        BaseURL: baseURL,
        HTTPClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        T: t,
    }
}

// SetAuthToken establece el token de autenticación
func (c *HTTPClient) SetAuthToken(token string) {
    c.AuthToken = token
}

// Request representa una petición HTTP
type Request struct {
    Method  string
    Path    string
    Body    interface{}
    Headers map[string]string
}

// Response representa una respuesta HTTP
type Response struct {
    StatusCode int
    Body       []byte
    Headers    http.Header
}

// Do ejecuta una petición HTTP y retorna la respuesta
func (c *HTTPClient) Do(req Request) (*Response, error) {
    url := c.BaseURL + req.Path

    var bodyReader io.Reader
    if req.Body != nil {
        jsonBody, err := json.Marshal(req.Body)
        if err != nil {
            return nil, fmt.Errorf("error marshaling request body: %w", err)
        }
        bodyReader = bytes.NewBuffer(jsonBody)
    }

    httpReq, err := http.NewRequest(req.Method, url, bodyReader)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    httpReq.Header.Set("Content-Type", "application/json")
    if c.AuthToken != "" {
        httpReq.Header.Set("Authorization", "Bearer "+c.AuthToken)
    }
    for key, value := range req.Headers {
        httpReq.Header.Set(key, value)
    }

    resp, err := c.HTTPClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("error executing request: %w", err)
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %w", err)
    }

    return &Response{
        StatusCode: resp.StatusCode,
        Body:       respBody,
        Headers:    resp.Header,
    }, nil
}

// GET ejecuta una petición GET
func (c *HTTPClient) GET(path string) (*Response, error) {
    return c.Do(Request{
        Method: http.MethodGet,
        Path:   path,
    })
}

// POST ejecuta una petición POST
func (c *HTTPClient) POST(path string, body interface{}) (*Response, error) {
    return c.Do(Request{
        Method: http.MethodPost,
        Path:   path,
        Body:   body,
    })
}

// PATCH ejecuta una petición PATCH
func (c *HTTPClient) PATCH(path string, body interface{}) (*Response, error) {
    return c.Do(Request{
        Method: http.MethodPatch,
        Path:   path,
        Body:   body,
    })
}

// DELETE ejecuta una petición DELETE
func (c *HTTPClient) DELETE(path string) (*Response, error) {
    return c.Do(Request{
        Method: http.MethodDelete,
        Path:   path,
    })
}

// ParseJSON deserializa el body de la respuesta a una estructura
func (r *Response) ParseJSON(target interface{}) error {
    return json.Unmarshal(r.Body, target)
}

// AssertStatus verifica que el status code sea el esperado
func (r *Response) AssertStatus(t *testing.T, expected int) {
    if r.StatusCode != expected {
        t.Errorf("Expected status %d, got %d. Body: %s", expected, r.StatusCode, string(r.Body))
    }
}

// AssertStatusOK verifica que el status sea 200
func (r *Response) AssertStatusOK(t *testing.T) {
    r.AssertStatus(t, http.StatusOK)
}

// AssertStatusCreated verifica que el status sea 201
func (r *Response) AssertStatusCreated(t *testing.T) {
    r.AssertStatus(t, http.StatusCreated)
}

// AssertStatusUnauthorized verifica que el status sea 401
func (r *Response) AssertStatusUnauthorized(t *testing.T) {
    r.AssertStatus(t, http.StatusUnauthorized)
}

// AssertStatusForbidden verifica que el status sea 403
func (r *Response) AssertStatusForbidden(t *testing.T) {
    r.AssertStatus(t, http.StatusForbidden)
}
