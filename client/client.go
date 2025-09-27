package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// LicenseChainClient represents the main client for the LicenseChain API
type LicenseChainClient struct {
	apiKey   string
	baseURL  string
	timeout  time.Duration
	retries  int
	client   *http.Client
}

// NewClient creates a new LicenseChain client
func NewClient(apiKey, baseURL string, timeout time.Duration, retries int) *LicenseChainClient {
	if baseURL == "" {
		baseURL = "https://api.licensechain.app"
	}
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	if retries <= 0 {
		retries = 3
	}

	return &LicenseChainClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		timeout: timeout,
		retries: retries,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// CreateClient creates a new client with default settings
func CreateClient(apiKey string, baseURL ...string) *LicenseChainClient {
	url := "https://api.licensechain.app"
	if len(baseURL) > 0 {
		url = baseURL[0]
	}
	return NewClient(apiKey, url, 30*time.Second, 3)
}

// FromEnvironment creates a client from environment variables
func FromEnvironment() *LicenseChainClient {
	apiKey := os.Getenv("LICENSECHAIN_API_KEY")
	baseURL := os.Getenv("LICENSECHAIN_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.licensechain.app"
	}
	return NewClient(apiKey, baseURL, 30*time.Second, 3)
}

// License Management

// CreateLicense creates a new license
func (c *LicenseChainClient) CreateLicense(req CreateLicenseRequest) (*License, error) {
	if err := ValidateNotEmpty(req.UserID, "user_id"); err != nil {
		return nil, err
	}
	if err := ValidateNotEmpty(req.ProductID, "product_id"); err != nil {
		return nil, err
	}

	req.Metadata = SanitizeMetadata(req.Metadata)
	
	var response struct {
		Data License `json:"data"`
	}
	
	err := c.makeRequest("POST", "/licenses", req, &response)
	if err != nil {
		return nil, err
	}
	
	return &response.Data, nil
}

// GetLicense retrieves a license by ID
func (c *LicenseChainClient) GetLicense(licenseID string) (*License, error) {
	if err := ValidateNotEmpty(licenseID, "license_id"); err != nil {
		return nil, err
	}
	if !ValidateUUID(licenseID) {
		return nil, NewValidationError("Invalid license_id format")
	}

	var response struct {
		Data License `json:"data"`
	}
	
	err := c.makeRequest("GET", "/licenses/"+licenseID, nil, &response)
	if err != nil {
		return nil, err
	}
	
	return &response.Data, nil
}

// ValidateLicense validates a license key
func (c *LicenseChainClient) ValidateLicense(licenseKey string) (bool, error) {
	if err := ValidateNotEmpty(licenseKey, "license_key"); err != nil {
		return false, err
	}

	req := map[string]string{"license_key": licenseKey}
	var response struct {
		Valid bool `json:"valid"`
	}
	
	err := c.makeRequest("POST", "/licenses/validate", req, &response)
	if err != nil {
		return false, err
	}
	
	return response.Valid, nil
}

// Health Check

// Ping pings the API
func (c *LicenseChainClient) Ping() (*PingResponse, error) {
	var response PingResponse
	err := c.makeRequest("GET", "/ping", nil, &response)
	if err != nil {
		return nil, err
	}
	
	return &response, nil
}

// Health checks the API health
func (c *LicenseChainClient) Health() (*HealthResponse, error) {
	var response HealthResponse
	err := c.makeRequest("GET", "/health", nil, &response)
	if err != nil {
		return nil, err
	}
	
	return &response, nil
}

// Private methods

func (c *LicenseChainClient) makeRequest(method, endpoint string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Version", "1.0")
	req.Header.Set("X-Platform", "go-sdk")
	req.Header.Set("User-Agent", "LicenseChain-Go-SDK/1.0.0")

	return RetryWithBackoff(func() error {
		resp, err := c.client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if result != nil {
				if resp.ContentLength == 0 {
					return nil
				}
				return json.NewDecoder(resp.Body).Decode(result)
			}
			return nil
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return NewHTTPError(resp.StatusCode, "Unknown error")
		}

		var errorResp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(bodyBytes, &errorResp); err != nil {
			return NewHTTPError(resp.StatusCode, string(bodyBytes))
		}

		switch resp.StatusCode {
		case 400:
			return NewValidationError(errorResp.Error)
		case 401, 403:
			return NewAuthenticationError(errorResp.Error)
		case 404:
			return NewNotFoundError(errorResp.Error)
		case 429:
			return NewRateLimitError(errorResp.Error)
		case 500, 502, 503, 504:
			return NewServerError(errorResp.Error)
		default:
			return NewHTTPError(resp.StatusCode, errorResp.Error)
		}
	}, c.retries, time.Second)
}