package client

import "time"

// License represents a license in the LicenseChain system
type License struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	ProductID string                 `json:"product_id"`
	LicenseKey string                `json:"license_key"`
	Status    string                 `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	ExpiresAt *time.Time             `json:"expires_at,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// CreateLicenseRequest represents a request to create a license
type CreateLicenseRequest struct {
	UserID    string                 `json:"user_id"`
	ProductID string                 `json:"product_id"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateLicenseRequest represents a request to update a license
type UpdateLicenseRequest struct {
	Status    string                 `json:"status,omitempty"`
	ExpiresAt *time.Time             `json:"expires_at,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// LicenseListResponse represents a paginated list of licenses
type LicenseListResponse struct {
	Data  []License `json:"data"`
	Total int       `json:"total"`
	Page  int       `json:"page"`
	Limit int       `json:"limit"`
}

// LicenseStats represents license statistics
type LicenseStats struct {
	Total   int     `json:"total"`
	Active  int     `json:"active"`
	Expired int     `json:"expired"`
	Revoked int     `json:"revoked"`
	Revenue float64 `json:"revenue"`
}

// User represents a user in the LicenseChain system
type User struct {
	ID        string                 `json:"id"`
	Email     string                 `json:"email"`
	Name      string                 `json:"name"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Email    string                 `json:"email"`
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	Email    string                 `json:"email,omitempty"`
	Name     string                 `json:"name,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// UserListResponse represents a paginated list of users
type UserListResponse struct {
	Data  []User `json:"data"`
	Total int    `json:"total"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

// UserStats represents user statistics
type UserStats struct {
	Total   int `json:"total"`
	Active  int `json:"active"`
	Inactive int `json:"inactive"`
}

// Product represents a product in the LicenseChain system
type Product struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Price       float64                `json:"price"`
	Currency    string                 `json:"currency"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// CreateProductRequest represents a request to create a product
type CreateProductRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Price       float64                `json:"price"`
	Currency    string                 `json:"currency"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateProductRequest represents a request to update a product
type UpdateProductRequest struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Price       float64                `json:"price,omitempty"`
	Currency    string                 `json:"currency,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ProductListResponse represents a paginated list of products
type ProductListResponse struct {
	Data  []Product `json:"data"`
	Total int       `json:"total"`
	Page  int       `json:"page"`
	Limit int       `json:"limit"`
}

// ProductStats represents product statistics
type ProductStats struct {
	Total   int     `json:"total"`
	Active  int     `json:"active"`
	Revenue float64 `json:"revenue"`
}

// Webhook represents a webhook in the LicenseChain system
type Webhook struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	Secret    string    `json:"secret,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateWebhookRequest represents a request to create a webhook
type CreateWebhookRequest struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
	Secret string   `json:"secret,omitempty"`
}

// UpdateWebhookRequest represents a request to update a webhook
type UpdateWebhookRequest struct {
	URL    string   `json:"url,omitempty"`
	Events []string `json:"events,omitempty"`
	Secret string   `json:"secret,omitempty"`
}

// WebhookListResponse represents a list of webhooks
type WebhookListResponse struct {
	Data  []Webhook `json:"data"`
	Total int       `json:"total"`
	Page  int       `json:"page"`
	Limit int       `json:"limit"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// PingResponse represents a ping response
type PingResponse struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}
