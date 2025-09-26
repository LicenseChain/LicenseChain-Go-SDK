// LicenseChain Go SDK - Basic Usage Example
package main

import (
    "fmt"
    "log"
    "github.com/LicenseChain/LicenseChain-Go-SDK"
)

func main() {
    fmt.Println("?? LicenseChain Go SDK - Basic Usage Example")
    fmt.Println("=" + strings.Repeat("=", 50))
    
    // Initialize the client
    client := licensechain.NewClient(&licensechain.Config{
        APIKey:  "your-api-key-here",
        AppName: "MyGoApp",
        Version: "1.0.0",
        Debug:   true,
    })
    
    // Connect to LicenseChain
    fmt.Println("\n Connecting to LicenseChain...")
    err := client.Connect()
    if err != nil {
        log.Fatalf(" Failed to connect to LicenseChain: %v", err)
    }
    
    fmt.Println(" Connected to LicenseChain successfully!")
    
    // Example 1: User Registration
    fmt.Println("\n Registering new user...")
    user, err := client.Register("testuser", "password123", "test@example.com")
    if err != nil {
        fmt.Printf(" Registration failed: %v\n", err)
    } else {
        fmt.Println(" User registered successfully!")
        fmt.Printf("Session ID: %s\n", user.SessionID)
    }
    
    // Example 2: License Validation
    fmt.Println("\n Validating license...")
    license, err := client.ValidateLicense("LICENSE-KEY-HERE")
    if err != nil {
        fmt.Printf(" License validation failed: %v\n", err)
    } else {
        fmt.Println(" License is valid!")
        fmt.Printf("License Key: %s\n", license.Key)
        fmt.Printf("Status: %s\n", license.Status)
    }
    
    // Cleanup
    fmt.Println("\n Cleaning up...")
    client.Logout()
    client.Disconnect()
    fmt.Println(" Cleanup completed!")
}
