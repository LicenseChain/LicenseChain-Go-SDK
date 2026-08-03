package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
		baseURL = "https://api.licensechain.app/v1"
	}
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	if retries <= 0 {
		retries = 3
	}

	return &LicenseChainClient{
		apiKey:  apiKey,
		baseURL: normalizeBaseURL(baseURL),
		timeout: timeout,
		retries: retries,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// CreateClient creates a new client with default settings
func CreateClient(apiKey string, baseURL ...string) *LicenseChainClient {
	url := "https://api.licensechain.app/v1"
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
		baseURL = "https://api.licensechain.app/v1"
	}
	return NewClient(apiKey, baseURL, 30*time.Second, 3)
}

// License Management

// CreateLicense creates a new license
func (c *LicenseChainClient) CreateLicense(req CreateLicenseRequest) (*License, error) {
	if err := ValidateNotEmpty(req.AppID, "app_id"); err != nil {
		return nil, err
	}
	if err := ValidateNotEmpty(req.UserEmail, "user_email"); err != nil {
		return nil, err
	}

	req.Metadata = SanitizeMetadata(req.Metadata)

	payload := map[string]interface{}{
		"appId":       req.AppID,
		"plan":        defaultString(req.Plan, "FREE"),
		"issuedEmail": req.UserEmail,
		"issuedTo":    req.UserName,
		"expiresAt":   req.ExpiresAt,
		"metadata":    req.Metadata,
	}

	var response map[string]interface{}
	err := c.makeRequest("POST", "/apps/"+req.AppID+"/licenses", payload, &response)
	if err != nil {
		return nil, err
	}

	licenseData := response
	if wrapped, ok := response["data"].(map[string]interface{}); ok {
		licenseData = wrapped
	}
	license := mapToLicense(licenseData)
	return &license, nil
}

// GetLicense retrieves a license by ID
func (c *LicenseChainClient) GetLicense(licenseID string) (*License, error) {
	if err := ValidateNotEmpty(licenseID, "license_id"); err != nil {
		return nil, err
	}

	var response map[string]interface{}
	
	err := c.makeRequest("GET", "/licenses/"+licenseID, nil, &response)
	if err != nil {
		return nil, err
	}

	licenseData := response
	if wrapped, ok := response["data"].(map[string]interface{}); ok {
		licenseData = wrapped
	}
	license := mapToLicense(licenseData)
	return &license, nil
}

// ValidateLicense validates a license key. Optional hwuid for ecosystem HMAC/HWUID contract.
func (c *LicenseChainClient) ValidateLicense(licenseKey string, hwuid ...string) (bool, error) {
	if err := ValidateNotEmpty(licenseKey, "license_key"); err != nil {
		return false, err
	}
	req := map[string]string{"key": licenseKey}
	if len(hwuid) > 0 && strings.TrimSpace(hwuid[0]) != "" {
		req["hwuid"] = strings.TrimSpace(hwuid[0])
	} else {
		req["hwuid"] = GenerateDefaultHWUID()
	}
	var response struct {
		Valid bool `json:"valid"`
	}
	err := c.makeRequest("POST", "/licenses/verify", req, &response)
	if err != nil {
		return false, err
	}
	return response.Valid, nil
}

// VerifyLicenseWithDetails returns the full JSON from POST /v1/licenses/verify (e.g. valid, license_token, license_jwks_uri).
func (c *LicenseChainClient) VerifyLicenseWithDetails(licenseKey string, hwuid ...string) (map[string]interface{}, error) {
	if err := ValidateNotEmpty(licenseKey, "license_key"); err != nil {
		return nil, err
	}
	req := map[string]string{"key": licenseKey}
	if len(hwuid) > 0 && strings.TrimSpace(hwuid[0]) != "" {
		req["hwuid"] = strings.TrimSpace(hwuid[0])
	} else {
		req["hwuid"] = GenerateDefaultHWUID()
	}
	var response map[string]interface{}
	err := c.makeRequest("POST", "/licenses/verify", req, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Health Check

// Ping pings the API
func (c *LicenseChainClient) Ping() (*PingResponse, error) {
	var health HealthResponse
	err := c.makeRequest("GET", "/health", nil, &health)
	if err != nil {
		return nil, err
	}
	
	return &PingResponse{
		Message: health.Status,
		Time:    health.Timestamp,
	}, nil
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
	normalizedEndpoint := normalizeEndpoint(c.baseURL, endpoint)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+normalizedEndpoint, reqBody)
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

func normalizeBaseURL(baseURL string) string {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return "https://api.licensechain.app/v1"
	}
	if strings.HasSuffix(baseURL, "/v1") {
		return baseURL
	}
	return baseURL + "/v1"
}

func normalizeEndpoint(baseURL, endpoint string) string {
	baseHasV1 := strings.HasSuffix(baseURL, "/v1")
	if strings.HasPrefix(endpoint, "/v1/") {
		if baseHasV1 {
			return strings.TrimPrefix(endpoint, "/v1")
		}
		return endpoint
	}
	if strings.HasPrefix(endpoint, "/") {
		if baseHasV1 {
			return endpoint
		}
		return "/v1" + endpoint
	}
	if baseHasV1 {
		return "/" + endpoint
	}
	return "/v1/" + endpoint
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func mapToLicense(data map[string]interface{}) License {
	license := License{
		ID:         toString(data["id"]),
		UserID:     toString(data["user_id"]),
		ProductID:  toString(data["product_id"]),
		LicenseKey: toString(firstNonNil(data["licenseKey"], data["license_key"], data["key"])),
		Status:     strings.ToLower(toString(data["status"])),
		Metadata:   mapField(data["metadata"]),
	}

	if expiresAt := toString(firstNonNil(data["expiresAt"], data["expires_at"])); expiresAt != "" {
		if t, err := time.Parse(time.RFC3339, expiresAt); err == nil {
			license.ExpiresAt = &t
		}
	}
	if createdAt := toString(firstNonNil(data["createdAt"], data["created_at"])); createdAt != "" {
		if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			license.CreatedAt = t
		}
	}
	if updatedAt := toString(firstNonNil(data["updatedAt"], data["updated_at"])); updatedAt != "" {
		if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
			license.UpdatedAt = t
		}
	}
	return license
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func mapField(v interface{}) map[string]interface{} {
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	return map[string]interface{}{}
}

func firstNonNil(values ...interface{}) interface{} {
	for _, v := range values {
		if v != nil {
			return v
		}
	}
	return nil
}