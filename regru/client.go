package regru

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client represents a client for the Reg.ru API.
type Client struct {
	username   string
	password   string
	baseURL    *url.URL
	HTTPClient *http.Client
}

// NewClient creates a new client for the Reg.ru API.
func NewClient(username, password, apiEndpoint, certFile, keyFile string) (*Client, error) {
	if apiEndpoint == "" {
		apiEndpoint = defaultApiEndpoint
	}

	baseURL, err := url.Parse(apiEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API endpoint: %w", err)
	}

	// Load client certificate and key
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate: %w", err)
	}

	// Configure TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Configure HTTP transport with TLS configuration
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	// Configure HTTP client with transport
	httpClient := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	client := &Client{
		username:   username,
		password:   password,
		baseURL:    baseURL,
		HTTPClient: httpClient,
	}

	return client, nil
}

func NewClientBasic(username, password, apiEndpoint string) (*Client, error) {
	if apiEndpoint == "" {
		apiEndpoint = defaultApiEndpoint
	}

	baseURL, err := url.Parse(apiEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API endpoint: %w", err)
	}

	client := &Client{
		username:   username,
		password:   password,
		baseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}

	return client, nil
}

func (c Client) doRequest(request any, fragments ...string) (*APIResponse, error) {
	// Generate the full URL
	endpoint := c.baseURL.JoinPath(fragments...)

	// Convert request to a map for parameter extraction
	requestMap := make(map[string]interface{})

	requestData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	err = json.Unmarshal(requestData, &requestMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	// Log the request for debugging
	fmt.Printf("Request to endpoint %s: %+v\n", endpoint.String(), requestMap)

	// Create a form for sending parameters
	form := url.Values{}

	// Add username and password
	form.Add("username", c.username)
	form.Add("password", c.password)

	// For complex structures, we use JSON via input_data
	hasComplexData := false
	for _, value := range requestMap {
		if _, ok := value.(map[string]interface{}); ok {
			hasComplexData = true
			break
		}
		if _, ok := value.([]interface{}); ok {
			hasComplexData = true
			break
		}
	}

	if hasComplexData {
		// Используем формат JSON для сложных структур
		form.Add("input_format", "json")
		jsonData, err := json.Marshal(requestMap)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request to JSON: %w", err)
		}
		form.Add("input_data", string(jsonData))
	} else {
		// Добавляем обычные параметры
		for key, value := range requestMap {
			// Пропускаем поля username и password
			if key == "Username" || key == "Password" || key == "username" || key == "password" {
				continue
			}

			switch v := value.(type) {
			case string:
				form.Add(key, v)
			case float64:
				form.Add(key, fmt.Sprintf("%v", v))
			case bool:
				form.Add(key, fmt.Sprintf("%v", v))
			case []interface{}:
				for _, item := range v {
					form.Add(key, fmt.Sprintf("%v", item))
				}
			default:
				form.Add(key, fmt.Sprintf("%v", v))
			}
		}
	}

	// Логируем отправляемую форму для отладки
	fmt.Printf("Form data: %+v\n", form)

	// Создаем HTTP запрос с формой
	req, err := http.NewRequest(http.MethodPost, endpoint.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	// Устанавливаем заголовок для формы
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Отправляем запрос
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Читаем тело ответа
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Логируем ответ для отладки
	fmt.Printf("Response: %s\n", string(raw))

	// Декодируем ответ
	var apiResp APIResponse
	err = json.Unmarshal(raw, &apiResp)
	if err != nil {
		return nil, err
	}

	return &apiResp, nil
}
