
package main

import (
    "fmt"
    "log"
    "licensechain" // assuming LicenseChain is the module name
)

func main() {
    client := licensechain.NewClient("https://licensechain.app/license/YOUR-LICENSE-HERE/validate")

    license := "test-license-key"
    isValid, err := client.ValidateLicense(license)
    if err != nil {
        log.Fatalf("License validation failed: %v", err)
    }

    if isValid {
        fmt.Println("License is valid!")
    } else {
        fmt.Println("License is invalid.")
    }
}
