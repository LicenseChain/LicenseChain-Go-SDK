package main

import (
	"fmt"
	"log"
	"time"

	"github.com/licensechain/LicenseChain-Go-SDK/client"
)

func main() {
	// Create a new client
	lc := client.CreateClient("your-api-key-here", "https://api.licensechain.app")

	// Test basic functionality
	fmt.Println("🚀 LicenseChain Go SDK - Basic Usage Example\n")

	// 1. Health Check
	fmt.Println("🏥 Health Check:")
	health, err := lc.Health()
	if err != nil {
		log.Printf("Health check failed: %v", err)
	} else {
		fmt.Printf("✅ API Status: %s\n", health.Status)
		fmt.Printf("   Version: %s\n", health.Version)
		fmt.Printf("   Timestamp: %s\n", health.Timestamp)
	}

	// 2. Ping
	fmt.Println("\n📡 Ping:")
	ping, err := lc.Ping()
	if err != nil {
		log.Printf("Ping failed: %v", err)
	} else {
		fmt.Printf("✅ Ping Response: %s\n", ping.Message)
		fmt.Printf("   Time: %s\n", ping.Time)
	}

	// 3. License Management
	fmt.Println("\n🔑 License Management:")

	// Create a license
	licenseReq := client.CreateLicenseRequest{
		UserID:    "user123",
		ProductID: "product456",
		Metadata: map[string]interface{}{
			"platform": "go",
			"version":  "1.0.0",
			"features": []string{"validation", "webhooks"},
		},
	}

	license, err := lc.CreateLicense(licenseReq)
	if err != nil {
		log.Printf("Failed to create license: %v", err)
	} else {
		fmt.Printf("✅ License created: %s\n", license.ID)
		fmt.Printf("   License Key: %s\n", license.LicenseKey)
		fmt.Printf("   Status: %s\n", license.Status)
		fmt.Printf("   Created: %s\n", license.CreatedAt.Format(time.RFC3339))
	}

	// Validate a license
	licenseKey := client.GenerateLicenseKey()
	fmt.Printf("\n🔍 Validating license key: %s\n", licenseKey)

	isValid, err := lc.ValidateLicense(licenseKey)
	if err != nil {
		log.Printf("Failed to validate license: %v", err)
	} else if isValid {
		fmt.Println("✅ License is valid")
	} else {
		fmt.Println("❌ License is invalid")
	}

	// 4. Utility Functions
	fmt.Println("\n🛠️ Utility Functions:")

	// Email validation
	email := "test@example.com"
	fmt.Printf("Email '%s' is valid: %t\n", email, client.ValidateEmail(email))

	// License key validation
	generatedKey := client.GenerateLicenseKey()
	fmt.Printf("License key '%s' is valid: %t\n", generatedKey, client.ValidateLicenseKey(generatedKey))

	// Generate UUID
	uuid := client.GenerateUUID()
	fmt.Printf("Generated UUID: %s\n", uuid)

	// Format bytes
	bytes := int64(1024 * 1024)
	fmt.Printf("%d bytes = %s\n", bytes, client.FormatBytes(bytes))

	// Format duration
	seconds := int64(3661)
	fmt.Printf("Duration: %s\n", client.FormatDuration(seconds))

	// String utilities
	text := "Hello World"
	fmt.Printf("Capitalize first: %s\n", client.CapitalizeFirst(text))
	fmt.Printf("To snake_case: %s\n", client.ToSnakeCase("HelloWorld"))
	fmt.Printf("To PascalCase: %s\n", client.ToPascalCase("hello_world"))
	fmt.Printf("Slugify: %s\n", client.Slugify("Hello World!"))

	// 5. Webhook Handling
	fmt.Println("\n🔄 Webhook Handling:")

	webhookHandler := client.NewWebhookHandler("webhook-secret", 300)

	// Simulate a webhook event
	webhookEvent := map[string]interface{}{
		"id":   "evt_123",
		"type": "license.created",
		"data": map[string]interface{}{
			"id":          "lic_123",
			"user_id":     "user_123",
			"product_id":  "prod_123",
			"license_key": "ABCDEFGHIJKLMNOPQRSTUVWXYZ012345",
			"status":      "active",
			"created_at":  time.Now().Format(time.RFC3339),
		},
		"timestamp": time.Now().Format(time.RFC3339),
		"signature": "signature_here",
	}

	err = webhookHandler.ProcessEvent(webhookEvent)
	if err != nil {
		log.Printf("Failed to process webhook event: %v", err)
	} else {
		fmt.Println("✅ Webhook event processed successfully")
	}

	// 6. Error Handling
	fmt.Println("\n🛡️ Error Handling:")

	// Test validation error
	err = client.ValidateNotEmpty("", "test_field")
	if err != nil {
		fmt.Printf("✅ Caught expected validation error: %v\n", err)
	}

	// Test positive validation
	err = client.ValidatePositive(-1, "test_amount")
	if err != nil {
		fmt.Printf("✅ Caught expected positive validation error: %v\n", err)
	}

	// Test range validation
	err = client.ValidateRange(5, 1, 3, "test_range")
	if err != nil {
		fmt.Printf("✅ Caught expected range validation error: %v\n", err)
	}

	// 7. JSON Utilities
	fmt.Println("\n📄 JSON Utilities:")

	// Test JSON serialization
	testData := map[string]interface{}{
		"name":    "Test",
		"value":   123,
		"active":  true,
		"items":   []string{"item1", "item2"},
	}

	jsonStr, err := client.JSONSerialize(testData)
	if err != nil {
		log.Printf("Failed to serialize JSON: %v", err)
	} else {
		fmt.Printf("✅ JSON serialized: %s\n", jsonStr)
	}

	// Test JSON validation
	validJSON := `{"name": "test", "value": 123}`
	invalidJSON := `{name: test, value: 123}`

	fmt.Printf("Valid JSON '%s': %t\n", validJSON, client.IsValidJSON(validJSON))
	fmt.Printf("Invalid JSON '%s': %t\n", invalidJSON, client.IsValidJSON(invalidJSON))

	// 8. URL Utilities
	fmt.Println("\n🌐 URL Utilities:")

	validURL := "https://api.licensechain.app"
	invalidURL := "not-a-url"

	fmt.Printf("Valid URL '%s': %t\n", validURL, client.IsValidURL(validURL))
	fmt.Printf("Invalid URL '%s': %t\n", invalidURL, client.IsValidURL(invalidURL))

	// Test URL encoding
	textToEncode := "Hello World! How are you?"
	encoded := client.URLEncode(textToEncode)
	fmt.Printf("URL encoded '%s': %s\n", textToEncode, encoded)

	decoded, err := client.URLDecode(encoded)
	if err != nil {
		log.Printf("Failed to decode URL: %v", err)
	} else {
		fmt.Printf("URL decoded '%s': %s\n", encoded, decoded)
	}

	// 9. Crypto Utilities
	fmt.Println("\n🔐 Crypto Utilities:")

	testString := "Hello, World!"
	fmt.Printf("SHA256 of '%s': %s\n", testString, client.SHA256(testString))
	fmt.Printf("SHA1 of '%s': %s\n", testString, client.SHA1(testString))
	fmt.Printf("MD5 of '%s': %s\n", testString, client.MD5(testString))

	// Test webhook signature
	payload := `{"test": "data"}`
	secret := "test-secret"
	signature := client.CreateWebhookSignature(payload, secret)
	fmt.Printf("Webhook signature: %s\n", signature)

	isValidSig := client.VerifyWebhookSignature(payload, signature, secret)
	fmt.Printf("Signature verification: %t\n", isValidSig)

	// 10. Base64 Utilities
	fmt.Println("\n📦 Base64 Utilities:")

	originalText := "Hello, World!"
	encodedB64 := client.Base64Encode(originalText)
	fmt.Printf("Base64 encoded '%s': %s\n", originalText, encodedB64)

	decodedB64, err := client.Base64Decode(encodedB64)
	if err != nil {
		log.Printf("Failed to decode Base64: %v", err)
	} else {
		fmt.Printf("Base64 decoded '%s': %s\n", encodedB64, decodedB64)
	}

	fmt.Println("\n✅ Basic usage example completed successfully!")
}