package client

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// ValidateEmail validates an email address
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateLicenseKey validates a license key format
func ValidateLicenseKey(licenseKey string) bool {
	if len(licenseKey) != 32 {
		return false
	}
	licenseKeyRegex := regexp.MustCompile(`^[A-Z0-9]+$`)
	return licenseKeyRegex.MatchString(licenseKey)
}

// ValidateUUID validates a UUID format
func ValidateUUID(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(strings.ToLower(uuid))
}

// ValidateAmount validates an amount
func ValidateAmount(amount float64) bool {
	return amount > 0 && !math.IsNaN(amount) && !math.IsInf(amount, 0)
}

// ValidateCurrency validates a currency code
func ValidateCurrency(currency string) bool {
	validCurrencies := []string{"USD", "EUR", "GBP", "CAD", "AUD", "JPY", "CHF", "CNY"}
	currency = strings.ToUpper(currency)
	for _, valid := range validCurrencies {
		if currency == valid {
			return true
		}
	}
	return false
}

// SanitizeInput sanitizes user input
func SanitizeInput(input string) string {
	input = strings.ReplaceAll(input, "&", "&amp;")
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#x27;")
	return input
}

// SanitizeMetadata sanitizes metadata
func SanitizeMetadata(metadata map[string]interface{}) map[string]interface{} {
	sanitized := make(map[string]interface{})
	for key, value := range metadata {
		switch v := value.(type) {
		case string:
			sanitized[key] = SanitizeInput(v)
		case []interface{}:
			sanitizedArray := make([]interface{}, len(v))
			for i, item := range v {
				if str, ok := item.(string); ok {
					sanitizedArray[i] = SanitizeInput(str)
				} else {
					sanitizedArray[i] = item
				}
			}
			sanitized[key] = sanitizedArray
		case map[string]interface{}:
			sanitized[key] = SanitizeMetadata(v)
		default:
			sanitized[key] = value
		}
	}
	return sanitized
}

// GenerateLicenseKey generates a random license key
func GenerateLicenseKey() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 32)
	for i := range result {
		result[i] = chars[time.Now().UnixNano()%int64(len(chars))]
	}
	return string(result)
}

// GenerateUUID generates a random UUID
func GenerateUUID() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		time.Now().UnixNano(),
		time.Now().UnixNano()>>32,
		time.Now().UnixNano()>>16,
		time.Now().UnixNano()>>8,
		time.Now().UnixNano())
}

// FormatTimestamp formats a timestamp
func FormatTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(time.RFC3339)
}

// ParseTimestamp parses a timestamp string
func ParseTimestamp(timestamp string) (int64, error) {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// ValidatePagination validates pagination parameters
func ValidatePagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}

// ValidateDateRange validates a date range
func ValidateDateRange(startDate, endDate string) error {
	start, err := ParseTimestamp(startDate)
	if err != nil {
		return fmt.Errorf("invalid start date: %v", err)
	}
	end, err := ParseTimestamp(endDate)
	if err != nil {
		return fmt.Errorf("invalid end date: %v", err)
	}
	if start > end {
		return fmt.Errorf("start date must be before or equal to end date")
	}
	return nil
}

// CreateWebhookSignature creates a webhook signature
func CreateWebhookSignature(payload, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyWebhookSignature verifies a webhook signature
func VerifyWebhookSignature(payload, signature, secret string) bool {
	expectedSignature := CreateWebhookSignature(payload, secret)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// RetryWithBackoff retries a function with exponential backoff
func RetryWithBackoff(fn func() error, maxRetries int, initialDelay time.Duration) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := fn(); err != nil {
			lastErr = err
			if i < maxRetries-1 {
				delay := time.Duration(float64(initialDelay) * math.Pow(2, float64(i)))
				time.Sleep(delay)
			}
		} else {
			return nil
		}
	}
	return lastErr
}

// FormatBytes formats bytes into human readable format
func FormatBytes(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	size := float64(bytes)
	unitIndex := 0
	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}
	return fmt.Sprintf("%.1f %s", size, units[unitIndex])
}

// FormatDuration formats duration in seconds
func FormatDuration(seconds int64) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	} else if seconds < 3600 {
		minutes := seconds / 60
		remainingSeconds := seconds % 60
		return fmt.Sprintf("%dm %ds", minutes, remainingSeconds)
	} else if seconds < 86400 {
		hours := seconds / 3600
		minutes := (seconds % 3600) / 60
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		days := seconds / 86400
		hours := (seconds % 86400) / 3600
		return fmt.Sprintf("%dd %dh", days, hours)
	}
}

// CapitalizeFirst capitalizes the first letter of a string
func CapitalizeFirst(text string) string {
	if len(text) == 0 {
		return text
	}
	return strings.ToUpper(text[:1]) + strings.ToLower(text[1:])
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(text string) string {
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	return strings.ToLower(re.ReplaceAllString(text, "${1}_${2}"))
}

// ToPascalCase converts a string to PascalCase
func ToPascalCase(text string) string {
	words := strings.Split(text, "_")
	for i, word := range words {
		words[i] = CapitalizeFirst(word)
	}
	return strings.Join(words, "")
}

// TruncateString truncates a string to a maximum length
func TruncateString(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength-3] + "..."
}

// Slugify converts a string to a URL-friendly slug
func Slugify(text string) string {
	text = strings.ToLower(text)
	re := regexp.MustCompile(`[^a-z0-9\s-]`)
	text = re.ReplaceAllString(text, "")
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, "-")
	re = regexp.MustCompile(`-+`)
	text = re.ReplaceAllString(text, "-")
	return strings.Trim(text, "-")
}

// ValidateNotEmpty validates that a string is not empty
func ValidateNotEmpty(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return NewValidationError(fmt.Sprintf("%s cannot be empty", fieldName))
	}
	return nil
}

// ValidatePositive validates that a number is positive
func ValidatePositive(value float64, fieldName string) error {
	if value <= 0 {
		return NewValidationError(fmt.Sprintf("%s must be positive", fieldName))
	}
	return nil
}

// ValidateRange validates that a number is within a range
func ValidateRange(value, min, max float64, fieldName string) error {
	if value < min || value > max {
		return NewValidationError(fmt.Sprintf("%s must be between %v and %v", fieldName, min, max))
	}
	return nil
}

// JSONSerialize serializes an object to JSON
func JSONSerialize(obj interface{}) (string, error) {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// JSONDeserialize deserializes JSON to an object
func JSONDeserialize(data string, obj interface{}) error {
	return json.Unmarshal([]byte(data), obj)
}

// IsValidJSON checks if a string is valid JSON
func IsValidJSON(jsonString string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(jsonString), &js) == nil
}

// IsValidURL checks if a string is a valid URL
func IsValidURL(urlString string) bool {
	_, err := url.Parse(urlString)
	return err == nil
}

// URLEncode encodes a string for URL use
func URLEncode(s string) string {
	return url.QueryEscape(s)
}

// URLDecode decodes a URL-encoded string
func URLDecode(s string) (string, error) {
	return url.QueryUnescape(s)
}

// GetCurrentTimestamp returns the current timestamp
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// GetCurrentDate returns the current date as ISO string
func GetCurrentDate() string {
	return time.Now().Format(time.RFC3339)
}

// Sleep pauses execution for the specified duration
func Sleep(duration time.Duration) {
	time.Sleep(duration)
}

// SHA256 computes SHA256 hash
func SHA256(data string) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

// SHA1 computes SHA1 hash
func SHA1(data string) string {
	h := sha1.Sum([]byte(data))
	return hex.EncodeToString(h[:])
}

// MD5 computes MD5 hash
func MD5(data string) string {
	h := md5.Sum([]byte(data))
	return hex.EncodeToString(h[:])
}

// Base64Encode encodes a string to base64
func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Base64Decode decodes a base64 string
func Base64Decode(data string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
