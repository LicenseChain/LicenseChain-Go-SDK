# LicenseChain Go SDK

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/LicenseChain/LicenseChain-Go-SDK)](https://goreportcard.com/report/github.com/LicenseChain/LicenseChain-Go-SDK)

Official Go SDK for LicenseChain - Secure license management for Go applications.

## 🚀 Features

- **🔐 Secure Authentication** - User registration, login, and session management
- **📜 License Management** - Create, validate, update, and revoke licenses
- **🛡️ Hardware ID Validation** - Prevent license sharing and unauthorized access
- **🔔 Webhook Support** - Real-time license events and notifications
- **📊 Analytics Integration** - Track license usage and performance metrics
- **⚡ High Performance** - Optimized for production workloads
- **🔄 Async Operations** - Non-blocking HTTP requests and data processing
- **🛠️ Easy Integration** - Simple API with comprehensive documentation

## 📦 Installation

### Method 1: Go Modules (Recommended)

```bash
# Add to your go.mod
go get github.com/LicenseChain/LicenseChain-Go-SDK

# Or with specific version
go get github.com/LicenseChain/LicenseChain-Go-SDK@v1.0.0
```

### Method 2: Manual Installation

```bash
# Clone the repository
git clone https://github.com/LicenseChain/LicenseChain-Go-SDK.git
cd LicenseChain-Go-SDK

# Install dependencies
go mod tidy

# Build the library
go build ./...
```

### Method 3: Vendor Directory

```bash
# Add to vendor directory
go mod vendor
```

## 🚀 Quick Start

### Basic Setup

```go
package main

import (
    "fmt"
    "log"
    "github.com/LicenseChain/LicenseChain-Go-SDK"
)

func main() {
    // Initialize the client
    client := licensechain.NewClient(&licensechain.Config{
        APIKey:  "your-api-key",
        AppName: "your-app-name",
        Version: "1.0.0",
    })
    
    // Connect to LicenseChain
    if err := client.Connect(); err != nil {
        log.Fatalf("Failed to connect to LicenseChain: %v", err)
    }
    
    fmt.Println("Connected to LicenseChain successfully!")
}
```

### User Authentication

```go
// Register a new user
user, err := client.Register("username", "password", "email@example.com")
if err != nil {
    log.Printf("Registration failed: %v", err)
} else {
    fmt.Println("User registered successfully!")
    fmt.Printf("User ID: %s\n", user.ID)
}

// Login existing user
user, err = client.Login("username", "password")
if err != nil {
    log.Printf("Login failed: %v", err)
} else {
    fmt.Println("User logged in successfully!")
    fmt.Printf("Session ID: %s\n", user.SessionID)
}
```

### License Management

```go
// Validate a license
license, err := client.ValidateLicense("LICENSE-KEY-HERE")
if err != nil {
    log.Printf("License validation failed: %v", err)
} else {
    fmt.Println("License is valid!")
    fmt.Printf("License Key: %s\n", license.Key)
    fmt.Printf("Status: %s\n", license.Status)
    fmt.Printf("Expires: %s\n", license.Expires)
    fmt.Printf("Features: %v\n", license.Features)
    fmt.Printf("User: %s\n", license.User)
}

// Get user's licenses
licenses, err := client.GetUserLicenses()
if err != nil {
    log.Printf("Failed to get licenses: %v", err)
} else {
    fmt.Printf("Found %d licenses:\n", len(licenses))
    for i, license := range licenses {
        fmt.Printf("  %d. %s - %s (Expires: %s)\n", 
            i+1, license.Key, license.Status, license.Expires)
    }
}
```

### Hardware ID Validation

```go
// Get hardware ID (automatically generated)
hardwareID := client.GetHardwareID()
fmt.Printf("Hardware ID: %s\n", hardwareID)

// Validate hardware ID with license
isValid, err := client.ValidateHardwareID("LICENSE-KEY-HERE", hardwareID)
if err != nil {
    log.Printf("Hardware ID validation failed: %v", err)
} else if isValid {
    fmt.Println("Hardware ID is valid for this license!")
} else {
    fmt.Println("Hardware ID is not valid for this license.")
}
```

### Webhook Integration

```go
// Set up webhook handler
client.SetWebhookHandler(func(event string, data map[string]string) {
    fmt.Printf("Webhook received: %s\n", event)
    
    switch event {
    case "license.created":
        fmt.Printf("New license created: %s\n", data["licenseKey"])
    case "license.updated":
        fmt.Printf("License updated: %s\n", data["licenseKey"])
    case "license.revoked":
        fmt.Printf("License revoked: %s\n", data["licenseKey"])
    }
})

// Start webhook listener
go client.StartWebhookListener()
```

## 📚 API Reference

### LicenseChain Client

#### Constructor

```go
client := licensechain.NewClient(&licensechain.Config{
    APIKey:  "your-api-key",
    AppName: "your-app-name",
    Version: "1.0.0",
    BaseURL: "https://api.licensechain.com", // Optional
})
```

#### Methods

##### Connection Management

```go
// Connect to LicenseChain
err := client.Connect()

// Disconnect from LicenseChain
client.Disconnect()

// Check connection status
isConnected := client.IsConnected()
```

##### User Authentication

```go
// Register a new user
user, err := client.Register(username, password, email)

// Login existing user
user, err := client.Login(username, password)

// Logout current user
client.Logout()

// Get current user info
user, err := client.GetCurrentUser()
```

##### License Management

```go
// Validate a license
license, err := client.ValidateLicense(licenseKey)

// Get user's licenses
licenses, err := client.GetUserLicenses()

// Create a new license
license, err := client.CreateLicense(userID, features, expires)

// Update a license
license, err := client.UpdateLicense(licenseKey, updates)

// Revoke a license
err := client.RevokeLicense(licenseKey)

// Extend a license
license, err := client.ExtendLicense(licenseKey, days)
```

##### Hardware ID Management

```go
// Get hardware ID
hardwareID := client.GetHardwareID()

// Validate hardware ID
isValid, err := client.ValidateHardwareID(licenseKey, hardwareID)

// Bind hardware ID to license
err := client.BindHardwareID(licenseKey, hardwareID)
```

##### Webhook Management

```go
// Set webhook handler
client.SetWebhookHandler(handler)

// Start webhook listener
go client.StartWebhookListener()

// Stop webhook listener
client.StopWebhookListener()
```

##### Analytics

```go
// Track event
err := client.TrackEvent(eventName, properties)

// Get analytics data
analytics, err := client.GetAnalytics(timeRange)
```

## 🔧 Configuration

### Environment Variables

Set these in your environment or through your build process:

```bash
# Required
export LICENSECHAIN_API_KEY=your-api-key
export LICENSECHAIN_APP_NAME=your-app-name
export LICENSECHAIN_APP_VERSION=1.0.0

# Optional
export LICENSECHAIN_BASE_URL=https://api.licensechain.com
export LICENSECHAIN_DEBUG=true
```

### Advanced Configuration

```go
client := licensechain.NewClient(&licensechain.Config{
    APIKey:     "your-api-key",
    AppName:    "your-app-name",
    Version:    "1.0.0",
    BaseURL:    "https://api.licensechain.com",
    Timeout:    30 * time.Second, // Request timeout
    Retries:    3,                // Number of retry attempts
    Debug:      false,            // Enable debug logging
    UserAgent:  "MyApp/1.0.0",   // Custom user agent
})
```

## 🛡️ Security Features

### Hardware ID Protection

The SDK automatically generates and manages hardware IDs to prevent license sharing:

```go
// Hardware ID is automatically generated and stored
hardwareID := client.GetHardwareID()

// Validate against license
isValid, err := client.ValidateHardwareID(licenseKey, hardwareID)
```

### Secure Communication

- All API requests use HTTPS
- API keys are securely stored and transmitted
- Session tokens are automatically managed
- Webhook signatures are verified

### License Validation

- Real-time license validation
- Hardware ID binding
- Expiration checking
- Feature-based access control

## 📊 Analytics and Monitoring

### Event Tracking

```go
// Track custom events
err := client.TrackEvent("app.started", map[string]interface{}{
    "level":       1,
    "playerCount": 10,
})

// Track license events
err := client.TrackEvent("license.validated", map[string]interface{}{
    "licenseKey": "LICENSE-KEY",
    "features":   "premium,unlimited",
})
```

### Performance Monitoring

```go
// Get performance metrics
metrics, err := client.GetPerformanceMetrics()
if err != nil {
    log.Printf("Failed to get metrics: %v", err)
} else {
    fmt.Printf("API Response Time: %v\n", metrics.AverageResponseTime)
    fmt.Printf("Success Rate: %.2f%%\n", metrics.SuccessRate*100)
    fmt.Printf("Error Count: %d\n", metrics.ErrorCount)
}
```

## 🔄 Error Handling

### Custom Error Types

```go
license, err := client.ValidateLicense("invalid-key")
if err != nil {
    switch e := err.(type) {
    case *licensechain.InvalidLicenseError:
        log.Println("License key is invalid")
    case *licensechain.ExpiredLicenseError:
        log.Println("License has expired")
    case *licensechain.NetworkError:
        log.Println("Network connection failed")
    case *licensechain.LicenseChainError:
        log.Printf("LicenseChain error: %v", e)
    default:
        log.Printf("Unknown error: %v", err)
    }
}
```

### Retry Logic

```go
// Automatic retry for network errors
client := licensechain.NewClient(&licensechain.Config{
    APIKey:  "your-api-key",
    AppName: "your-app-name",
    Version: "1.0.0",
    Retries: 3,              // Retry up to 3 times
    Timeout: 30 * time.Second, // Wait 30 seconds for each request
})
```

## 🧪 Testing

### Unit Tests

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Integration Tests

```bash
# Test with real API
go test -tags=integration ./...

# Test specific package
go test -v ./client
```

## 📝 Examples

See the `examples/` directory for complete examples:

- `basic_usage.go` - Basic SDK usage
- `advanced_features.go` - Advanced features and configuration
- `webhook_integration.go` - Webhook handling

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

1. Clone the repository
2. Install Go 1.19 or later
3. Build: `go build ./...`
4. Test: `go test ./...`

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Documentation**: [https://docs.licensechain.com/go](https://docs.licensechain.com/go)
- **Issues**: [GitHub Issues](https://github.com/LicenseChain/LicenseChain-Go-SDK/issues)
- **Discord**: [LicenseChain Discord](https://discord.gg/licensechain)
- **Email**: support@licensechain.com

## 🔗 Related Projects

- [LicenseChain JavaScript SDK](https://github.com/LicenseChain/LicenseChain-JavaScript-SDK)
- [LicenseChain Python SDK](https://github.com/LicenseChain/LicenseChain-Python-SDK)
- [LicenseChain Node.js SDK](https://github.com/LicenseChain/LicenseChain-NodeJS-SDK)
- [LicenseChain Customer Panel](https://github.com/LicenseChain/LicenseChain-Customer-Panel)

---

**Made with ❤️ for the Go community**
