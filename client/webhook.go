package client

import (
	"encoding/json"
	"fmt"
	"time"
)

// WebhookHandler handles webhook events
type WebhookHandler struct {
	secret    string
	tolerance int64 // seconds
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(secret string, tolerance int64) *WebhookHandler {
	if tolerance <= 0 {
		tolerance = 300 // 5 minutes default
	}
	return &WebhookHandler{
		secret:    secret,
		tolerance: tolerance,
	}
}

// VerifySignature verifies a webhook signature
func (wh *WebhookHandler) VerifySignature(payload, signature string) bool {
	return VerifyWebhookSignature(payload, signature, wh.secret)
}

// VerifyTimestamp verifies a webhook timestamp
func (wh *WebhookHandler) VerifyTimestamp(timestamp string) error {
	webhookTime, err := ParseTimestamp(timestamp)
	if err != nil {
		return fmt.Errorf("invalid timestamp format: %v", err)
	}
	
	currentTime := GetCurrentTimestamp()
	timeDiff := currentTime - webhookTime
	if timeDiff < 0 {
		timeDiff = -timeDiff
	}
	
	if timeDiff > wh.tolerance {
		return fmt.Errorf("webhook timestamp too old: %d seconds", timeDiff)
	}
	
	return nil
}

// VerifyWebhook verifies a webhook
func (wh *WebhookHandler) VerifyWebhook(payload, signature, timestamp string) error {
	if err := wh.VerifyTimestamp(timestamp); err != nil {
		return err
	}
	
	if !wh.VerifySignature(payload, signature) {
		return NewAuthenticationError("Invalid webhook signature")
	}
	
	return nil
}

// ProcessEvent processes a webhook event
func (wh *WebhookHandler) ProcessEvent(eventData map[string]interface{}) error {
	payload, err := JSONSerialize(eventData["data"])
	if err != nil {
		return fmt.Errorf("invalid event data: %v", err)
	}
	
	signature, ok := eventData["signature"].(string)
	if !ok {
		return NewValidationError("Missing signature")
	}
	
	timestamp, ok := eventData["timestamp"].(string)
	if !ok {
		return NewValidationError("Missing timestamp")
	}
	
	if err := wh.VerifyWebhook(payload, signature, timestamp); err != nil {
		return err
	}
	
	eventType, ok := eventData["type"].(string)
	if !ok {
		return NewValidationError("Missing event type")
	}
	
	switch eventType {
	case "license.created":
		return wh.handleLicenseCreated(eventData)
	case "license.updated":
		return wh.handleLicenseUpdated(eventData)
	case "license.revoked":
		return wh.handleLicenseRevoked(eventData)
	case "license.expired":
		return wh.handleLicenseExpired(eventData)
	case "user.created":
		return wh.handleUserCreated(eventData)
	case "user.updated":
		return wh.handleUserUpdated(eventData)
	case "user.deleted":
		return wh.handleUserDeleted(eventData)
	case "product.created":
		return wh.handleProductCreated(eventData)
	case "product.updated":
		return wh.handleProductUpdated(eventData)
	case "product.deleted":
		return wh.handleProductDeleted(eventData)
	case "payment.completed":
		return wh.handlePaymentCompleted(eventData)
	case "payment.failed":
		return wh.handlePaymentFailed(eventData)
	case "payment.refunded":
		return wh.handlePaymentRefunded(eventData)
	default:
		fmt.Printf("Unknown webhook event type: %s\n", eventType)
		return nil
	}
}

// Event handlers
func (wh *WebhookHandler) handleLicenseCreated(eventData map[string]interface{}) error {
	id, _ := eventData["id"].(string)
	fmt.Printf("License created: %s\n", id)
	// Add custom logic for license created event
	return nil
}

func (wh *WebhookHandler) handleLicenseUpdated(eventData map[string]interface{}) error {
	id, _ := eventData["id"].(string)
	fmt.Printf("License updated: %s\n", id)
	// Add custom logic for license updated event
	return nil
}

func (wh *WebhookHandler) handleLicenseRevoked(eventData map[string]interface{}) error {
	id, _ := eventData["id"].(string)
	fmt.Printf("License revoked: %s\n", id)
	// Add custom logic for license revoked event
	return nil
}

func (wh *WebhookHandler) handleLicenseExpired(eventData map[string]interface{}) error {
	id, _ := eventData["id"].(string)
	fmt.Printf("License expired: %s\n", id)
	// Add custom logic for license expired event
	return nil
}

func (wh *WebhookHandler) handleUserCreated(eventData map[string]interface{}) error {
	id, _ := eventData["id"].(string)
	fmt.Printf("User created: %s\n", id)
	// Add custom logic for user created event
	return nil
}

func (wh *WebhookHandler) handleUserUpdated(eventData map[string]interface{}) error {
	id, _ := eventData["id"].(string)
	fmt.Printf("User updated: %s\n", id)
	// Add custom logic for user updated event
	return nil
}

func (wh *WebhookHandler) handleUserDeleted(eventData map[string]interface{}) error {
	id, _ := eventData["id"].(string)
	fmt.Printf("User deleted: %s\n", id)
	// Add custom logic for user deleted event
	return nil
}

func (wh *WebhookHandler) handleProductCreated(eventData map[string]interface{}) error {
	id, _ := eventData["id"].(string)
	fmt.Printf("Product created: %s\n", id)
	// Add custom logic for product created event
	return nil
}

func (wh *WebhookHandler) handleProductUpdated(eventData map[string]interface{}) error {
	id, _ := eventData["id"].(string)
	fmt.Printf("Product updated: %s\n", id)
	// Add custom logic for product updated event
	return nil
}

func (wh *WebhookHandler) handleProductDeleted(eventData map[string]interface{}) error {
	id, _ := eventData["id"].(string)
	fmt.Printf("Product deleted: %s\n", id)
	// Add custom logic for product deleted event
	return nil
}

func (wh *WebhookHandler) handlePaymentCompleted(eventData map[string]interface{}) error {
	id, _ := eventData["id"].(string)
	fmt.Printf("Payment completed: %s\n", id)
	// Add custom logic for payment completed event
	return nil
}

func (wh *WebhookHandler) handlePaymentFailed(eventData map[string]interface{}) error {
	id, _ := eventData["id"].(string)
	fmt.Printf("Payment failed: %s\n", id)
	// Add custom logic for payment failed event
	return nil
}

func (wh *WebhookHandler) handlePaymentRefunded(eventData map[string]interface{}) error {
	id, _ := eventData["id"].(string)
	fmt.Printf("Payment refunded: %s\n", id)
	// Add custom logic for payment refunded event
	return nil
}

// WebhookEvents contains webhook event type constants
var WebhookEvents = struct {
	LicenseCreated   string
	LicenseUpdated   string
	LicenseRevoked   string
	LicenseExpired   string
	UserCreated      string
	UserUpdated      string
	UserDeleted      string
	ProductCreated   string
	ProductUpdated   string
	ProductDeleted   string
	PaymentCompleted string
	PaymentFailed    string
	PaymentRefunded  string
}{
	LicenseCreated:   "license.created",
	LicenseUpdated:   "license.updated",
	LicenseRevoked:   "license.revoked",
	LicenseExpired:   "license.expired",
	UserCreated:      "user.created",
	UserUpdated:      "user.updated",
	UserDeleted:      "user.deleted",
	ProductCreated:   "product.created",
	ProductUpdated:   "product.updated",
	ProductDeleted:   "product.deleted",
	PaymentCompleted: "payment.completed",
	PaymentFailed:    "payment.failed",
	PaymentRefunded:  "payment.refunded",
}
