package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// LicenseChainClient represents the main client for interacting with the LicenseChain API
type LicenseChainClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	timeout    time.Duration
	retries    int
	retryDelay time.Duration
}

// Config holds the configuration for the LicenseChain client
type Config struct {
	APIKey     string
	BaseURL    string
	Timeout    time.Duration
	Retries    int
	RetryDelay time.Duration
}

// NewClient creates a new LicenseChain client with the given configuration
func NewClient(config Config) (*LicenseChainClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	if config.BaseURL == "" {
		config.BaseURL = "https://api.licensechain.app"
	}

	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	if config.Retries == 0 {
		config.Retries = 3
	}

	if config.RetryDelay == 0 {
		config.RetryDelay = 1 * time.Second
	}

	client := &LicenseChainClient{
		baseURL: strings.TrimSuffix(config.BaseURL, "/"),
		apiKey:  config.APIKey,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		timeout:    config.Timeout,
		retries:    config.Retries,
		retryDelay: config.RetryDelay,
	}

	return client, nil
}

// User represents a user in the LicenseChain system
type User struct {
	ID            string    `json:"id,omitempty"`
	Email         string    `json:"email"`
	Name          string    `json:"name,omitempty"`
	Company       string    `json:"company,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	EmailVerified bool      `json:"email_verified,omitempty"`
	Status        string    `json:"status,omitempty"`
}

// Application represents an application in the LicenseChain system
type Application struct {
	ID              string    `json:"id,omitempty"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	APIKey          string    `json:"api_key,omitempty"`
	WebhookURL      string    `json:"webhook_url,omitempty"`
	AllowedOrigins  []string  `json:"allowed_origins,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
	Status          string    `json:"status,omitempty"`
	LicenseCount    int       `json:"license_count,omitempty"`
}

// License represents a license in the LicenseChain system
type License struct {
	ID          string                 `json:"id,omitempty"`
	Key         string                 `json:"key,omitempty"`
	AppID       string                 `json:"app_id"`
	UserID      string                 `json:"user_id,omitempty"`
	UserEmail   string                 `json:"user_email"`
	UserName    string                 `json:"user_name,omitempty"`
	Status      string                 `json:"status,omitempty"`
	ExpiresAt   time.Time              `json:"expires_at,omitempty"`
	CreatedAt   time.Time              `json:"created_at,omitempty"`
	UpdatedAt   time.Time              `json:"updated_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Features    []string               `json:"features,omitempty"`
	UsageCount  int                    `json:"usage_count,omitempty"`
}

// Webhook represents a webhook in the LicenseChain system
type Webhook struct {
	ID             string    `json:"id,omitempty"`
	AppID          string    `json:"app_id"`
	URL            string    `json:"url"`
	Events         []string  `json:"events,omitempty"`
	Secret         string    `json:"secret,omitempty"`
	Status         string    `json:"status,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	LastTriggeredAt time.Time `json:"last_triggered_at,omitempty"`
	FailureCount   int       `json:"failure_count,omitempty"`
}

// Analytics represents analytics data
type Analytics struct {
	TotalLicenses        int                    `json:"total_licenses,omitempty"`
	ActiveLicenses       int                    `json:"active_licenses,omitempty"`
	ExpiredLicenses      int                    `json:"expired_licenses,omitempty"`
	RevokedLicenses      int                    `json:"revoked_licenses,omitempty"`
	ValidationsToday     int                    `json:"validations_today,omitempty"`
	ValidationsThisWeek  int                    `json:"validations_this_week,omitempty"`
	ValidationsThisMonth int                    `json:"validations_this_month,omitempty"`
	TopFeatures          []string               `json:"top_features,omitempty"`
	UsageByDay           []map[string]interface{} `json:"usage_by_day,omitempty"`
	UsageByWeek          []map[string]interface{} `json:"usage_by_week,omitempty"`
	UsageByMonth         []map[string]interface{} `json:"usage_by_month,omitempty"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int   `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// LicenseValidationResult represents the result of license validation
type LicenseValidationResult struct {
	Valid     bool                   `json:"valid"`
	License   map[string]interface{} `json:"license,omitempty"`
	User      map[string]interface{} `json:"user,omitempty"`
	App       map[string]interface{} `json:"app,omitempty"`
	ExpiresAt string                 `json:"expires_at,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error     string                 `json:"error"`
	Message   string                 `json:"message,omitempty"`
	Code      string                 `json:"code,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp string                 `json:"timestamp,omitempty"`
}

// LicenseChainError represents a custom error from the LicenseChain API
type LicenseChainError struct {
	Message    string
	StatusCode int
	Code       string
	Details    map[string]interface{}
}

func (e *LicenseChainError) Error() string {
	return e.Message
}

// HTTP Methods

func (c *LicenseChainClient) get(ctx context.Context, path string, params map[string]string, result interface{}) error {
	req, err := c.newRequest(ctx, "GET", path, nil, params)
	if err != nil {
		return err
	}
	return c.do(req, result)
}

func (c *LicenseChainClient) post(ctx context.Context, path string, body interface{}, result interface{}) error {
	req, err := c.newRequest(ctx, "POST", path, body, nil)
	if err != nil {
		return err
	}
	return c.do(req, result)
}

func (c *LicenseChainClient) patch(ctx context.Context, path string, body interface{}, result interface{}) error {
	req, err := c.newRequest(ctx, "PATCH", path, body, nil)
	if err != nil {
		return err
	}
	return c.do(req, result)
}

func (c *LicenseChainClient) delete(ctx context.Context, path string) error {
	req, err := c.newRequest(ctx, "DELETE", path, nil, nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

func (c *LicenseChainClient) newRequest(ctx context.Context, method, path string, body interface{}, params map[string]string) (*http.Request, error) {
	url := c.baseURL + path

	// Add query parameters
	if params != nil && len(params) > 0 {
		u, err := url.Parse(url)
		if err != nil {
			return nil, err
		}
		q := u.Query()
		for key, value := range params {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
		url = u.String()
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", "LicenseChain-Go-SDK/1.0.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *LicenseChainClient) do(req *http.Request, result interface{}) error {
	var lastErr error

	for i := 0; i <= c.retries; i++ {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if i < c.retries {
				time.Sleep(c.retryDelay * time.Duration(1<<i)) // Exponential backoff
				continue
			}
			return err
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			if i < c.retries && isRetryableError(resp.StatusCode) {
				time.Sleep(c.retryDelay * time.Duration(1<<i))
				continue
			}
			return err
		}

		if !resp.IsSuccess() {
			var errorResp ErrorResponse
			if err := json.Unmarshal(body, &errorResp); err != nil {
				lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
			} else {
				lastErr = &LicenseChainError{
					Message:    errorResp.Error,
					StatusCode: resp.StatusCode,
					Code:       errorResp.Code,
					Details:    errorResp.Details,
				}
			}

			if i < c.retries && isRetryableError(resp.StatusCode) {
				time.Sleep(c.retryDelay * time.Duration(1<<i))
				continue
			}
			return lastErr
		}

		if result != nil && len(body) > 0 {
			if err := json.Unmarshal(body, result); err != nil {
				return err
			}
		}

		return nil
	}

	return lastErr
}

func isRetryableError(statusCode int) bool {
	return statusCode == 429 || statusCode >= 500
}

// Authentication Methods

// RegisterUser registers a new user
func (c *LicenseChainClient) RegisterUser(ctx context.Context, req UserRegistrationRequest) (*User, error) {
	var user User
	err := c.post(ctx, "/auth/register", req, &user)
	return &user, err
}

// Login authenticates a user
func (c *LicenseChainClient) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	var resp LoginResponse
	err := c.post(ctx, "/auth/login", req, &resp)
	return &resp, err
}

// Logout logs out the current user
func (c *LicenseChainClient) Logout(ctx context.Context) error {
	return c.post(ctx, "/auth/logout", nil, nil)
}

// RefreshToken refreshes the authentication token
func (c *LicenseChainClient) RefreshToken(ctx context.Context, refreshToken string) (*TokenRefreshResponse, error) {
	req := map[string]string{"refresh_token": refreshToken}
	var resp TokenRefreshResponse
	err := c.post(ctx, "/auth/refresh", req, &resp)
	return &resp, err
}

// GetUserProfile gets the current user profile
func (c *LicenseChainClient) GetUserProfile(ctx context.Context) (*User, error) {
	var user User
	err := c.get(ctx, "/auth/me", nil, &user)
	return &user, err
}

// UpdateUserProfile updates the user profile
func (c *LicenseChainClient) UpdateUserProfile(ctx context.Context, req UserUpdateRequest) (*User, error) {
	var user User
	err := c.patch(ctx, "/auth/me", req, &user)
	return &user, err
}

// ChangePassword changes the user password
func (c *LicenseChainClient) ChangePassword(ctx context.Context, req PasswordChangeRequest) error {
	return c.patch(ctx, "/auth/password", req, nil)
}

// RequestPasswordReset requests a password reset
func (c *LicenseChainClient) RequestPasswordReset(ctx context.Context, email string) error {
	req := map[string]string{"email": email}
	return c.post(ctx, "/auth/forgot-password", req, nil)
}

// ResetPassword resets the password with a token
func (c *LicenseChainClient) ResetPassword(ctx context.Context, req PasswordResetRequest) error {
	return c.post(ctx, "/auth/reset-password", req, nil)
}

// Application Management

// CreateApplication creates a new application
func (c *LicenseChainClient) CreateApplication(ctx context.Context, req ApplicationCreateRequest) (*Application, error) {
	var app Application
	err := c.post(ctx, "/apps", req, &app)
	return &app, err
}

// ListApplications lists applications with pagination
func (c *LicenseChainClient) ListApplications(ctx context.Context, req ApplicationListRequest) (*PaginatedResponse[Application], error) {
	params := make(map[string]string)
	if req.Page > 0 {
		params["page"] = fmt.Sprintf("%d", req.Page)
	}
	if req.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", req.Limit)
	}
	if req.Status != "" {
		params["status"] = req.Status
	}
	if req.SortBy != "" {
		params["sort_by"] = req.SortBy
	}
	if req.SortOrder != "" {
		params["sort_order"] = req.SortOrder
	}

	var resp PaginatedResponse[Application]
	err := c.get(ctx, "/apps", params, &resp)
	return &resp, err
}

// GetApplication gets application details
func (c *LicenseChainClient) GetApplication(ctx context.Context, appID string) (*Application, error) {
	var app Application
	err := c.get(ctx, "/apps/"+appID, nil, &app)
	return &app, err
}

// UpdateApplication updates an application
func (c *LicenseChainClient) UpdateApplication(ctx context.Context, appID string, req ApplicationUpdateRequest) (*Application, error) {
	var app Application
	err := c.patch(ctx, "/apps/"+appID, req, &app)
	return &app, err
}

// DeleteApplication deletes an application
func (c *LicenseChainClient) DeleteApplication(ctx context.Context, appID string) error {
	return c.delete(ctx, "/apps/"+appID)
}

// RegenerateAPIKey regenerates the API key for an application
func (c *LicenseChainClient) RegenerateAPIKey(ctx context.Context, appID string) (*APIKeyResponse, error) {
	var resp APIKeyResponse
	err := c.post(ctx, "/apps/"+appID+"/regenerate-key", nil, &resp)
	return &resp, err
}

// License Management

// CreateLicense creates a new license
func (c *LicenseChainClient) CreateLicense(ctx context.Context, req LicenseCreateRequest) (*License, error) {
	var license License
	err := c.post(ctx, "/licenses", req, &license)
	return &license, err
}

// ListLicenses lists licenses with filters
func (c *LicenseChainClient) ListLicenses(ctx context.Context, req LicenseListRequest) (*PaginatedResponse[License], error) {
	params := make(map[string]string)
	if req.Page > 0 {
		params["page"] = fmt.Sprintf("%d", req.Page)
	}
	if req.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", req.Limit)
	}
	if req.AppID != "" {
		params["app_id"] = req.AppID
	}
	if req.Status != "" {
		params["status"] = req.Status
	}
	if req.UserID != "" {
		params["user_id"] = req.UserID
	}
	if req.UserEmail != "" {
		params["user_email"] = req.UserEmail
	}
	if req.SortBy != "" {
		params["sort_by"] = req.SortBy
	}
	if req.SortOrder != "" {
		params["sort_order"] = req.SortOrder
	}

	var resp PaginatedResponse[License]
	err := c.get(ctx, "/licenses", params, &resp)
	return &resp, err
}

// GetLicense gets license details
func (c *LicenseChainClient) GetLicense(ctx context.Context, licenseID string) (*License, error) {
	var license License
	err := c.get(ctx, "/licenses/"+licenseID, nil, &license)
	return &license, err
}

// UpdateLicense updates a license
func (c *LicenseChainClient) UpdateLicense(ctx context.Context, licenseID string, req LicenseUpdateRequest) (*License, error) {
	var license License
	err := c.patch(ctx, "/licenses/"+licenseID, req, &license)
	return &license, err
}

// DeleteLicense deletes a license
func (c *LicenseChainClient) DeleteLicense(ctx context.Context, licenseID string) error {
	return c.delete(ctx, "/licenses/"+licenseID)
}

// ValidateLicense validates a license key
func (c *LicenseChainClient) ValidateLicense(ctx context.Context, licenseKey, appID string) (*LicenseValidationResult, error) {
	req := map[string]string{"license_key": licenseKey}
	if appID != "" {
		req["app_id"] = appID
	}

	var result LicenseValidationResult
	err := c.post(ctx, "/licenses/validate", req, &result)
	return &result, err
}

// RevokeLicense revokes a license
func (c *LicenseChainClient) RevokeLicense(ctx context.Context, licenseID, reason string) error {
	req := map[string]string{}
	if reason != "" {
		req["reason"] = reason
	}
	return c.patch(ctx, "/licenses/"+licenseID+"/revoke", req, nil)
}

// ActivateLicense activates a license
func (c *LicenseChainClient) ActivateLicense(ctx context.Context, licenseID string) error {
	return c.patch(ctx, "/licenses/"+licenseID+"/activate", nil, nil)
}

// ExtendLicense extends license expiration
func (c *LicenseChainClient) ExtendLicense(ctx context.Context, licenseID, expiresAt string) error {
	req := map[string]string{"expires_at": expiresAt}
	return c.patch(ctx, "/licenses/"+licenseID+"/extend", req, nil)
}

// Webhook Management

// CreateWebhook creates a webhook
func (c *LicenseChainClient) CreateWebhook(ctx context.Context, req WebhookCreateRequest) (*Webhook, error) {
	var webhook Webhook
	err := c.post(ctx, "/webhooks", req, &webhook)
	return &webhook, err
}

// ListWebhooks lists webhooks
func (c *LicenseChainClient) ListWebhooks(ctx context.Context, req WebhookListRequest) (*PaginatedResponse[Webhook], error) {
	params := make(map[string]string)
	if req.Page > 0 {
		params["page"] = fmt.Sprintf("%d", req.Page)
	}
	if req.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", req.Limit)
	}
	if req.AppID != "" {
		params["app_id"] = req.AppID
	}
	if req.Status != "" {
		params["status"] = req.Status
	}

	var resp PaginatedResponse[Webhook]
	err := c.get(ctx, "/webhooks", params, &resp)
	return &resp, err
}

// GetWebhook gets webhook details
func (c *LicenseChainClient) GetWebhook(ctx context.Context, webhookID string) (*Webhook, error) {
	var webhook Webhook
	err := c.get(ctx, "/webhooks/"+webhookID, nil, &webhook)
	return &webhook, err
}

// UpdateWebhook updates a webhook
func (c *LicenseChainClient) UpdateWebhook(ctx context.Context, webhookID string, req WebhookUpdateRequest) (*Webhook, error) {
	var webhook Webhook
	err := c.patch(ctx, "/webhooks/"+webhookID, req, &webhook)
	return &webhook, err
}

// DeleteWebhook deletes a webhook
func (c *LicenseChainClient) DeleteWebhook(ctx context.Context, webhookID string) error {
	return c.delete(ctx, "/webhooks/"+webhookID)
}

// TestWebhook tests a webhook
func (c *LicenseChainClient) TestWebhook(ctx context.Context, webhookID string) error {
	return c.post(ctx, "/webhooks/"+webhookID+"/test", nil, nil)
}

// Analytics

// GetAnalytics gets analytics data
func (c *LicenseChainClient) GetAnalytics(ctx context.Context, req AnalyticsRequest) (*Analytics, error) {
	params := make(map[string]string)
	if req.AppID != "" {
		params["app_id"] = req.AppID
	}
	if req.StartDate != "" {
		params["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		params["end_date"] = req.EndDate
	}
	if req.Metric != "" {
		params["metric"] = req.Metric
	}
	if req.Period != "" {
		params["period"] = req.Period
	}

	var analytics Analytics
	err := c.get(ctx, "/analytics", params, &analytics)
	return &analytics, err
}

// GetLicenseAnalytics gets license analytics
func (c *LicenseChainClient) GetLicenseAnalytics(ctx context.Context, licenseID string) (*Analytics, error) {
	var analytics Analytics
	err := c.get(ctx, "/licenses/"+licenseID+"/analytics", nil, &analytics)
	return &analytics, err
}

// GetUsageStats gets usage statistics
func (c *LicenseChainClient) GetUsageStats(ctx context.Context, req UsageStatsRequest) (*UsageStats, error) {
	params := make(map[string]string)
	if req.AppID != "" {
		params["app_id"] = req.AppID
	}
	if req.Period != "" {
		params["period"] = req.Period
	}
	if req.Granularity != "" {
		params["granularity"] = req.Granularity
	}

	var stats UsageStats
	err := c.get(ctx, "/analytics/usage", params, &stats)
	return &stats, err
}

// System Status

// GetSystemStatus gets system status
func (c *LicenseChainClient) GetSystemStatus(ctx context.Context) (*SystemStatus, error) {
	var status SystemStatus
	err := c.get(ctx, "/status", nil, &status)
	return &status, err
}

// GetHealthCheck gets health check
func (c *LicenseChainClient) GetHealthCheck(ctx context.Context) (*HealthCheck, error) {
	var health HealthCheck
	err := c.get(ctx, "/health", nil, &health)
	return &health, err
}
