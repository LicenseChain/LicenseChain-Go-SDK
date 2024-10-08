
# LicenseChain Golang SDK

This Golang SDK allows developers to integrate LicenseChain license validation into their Go applications using LicenseChain's APIs.

## Features
- Validate licenses through LicenseChain API
- Easy integration with your Golang projects

## Installation

```bash
go get github.com/LicenseChain/LicenseChain-Golang-SDK
```

## Usage Example

```go
package main

import (
    "fmt"
    "log"
    "licensechain" // assuming LicenseChain is the module name
)

func main() {
    client := licensechain.NewClient("your-api-key")

    license := "license-key-to-validate"
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
```

# Bugs
If the default example that wasn’t included in your software isn’t working as expected, please pop over to https://t.me/LicenseChainBot and lodge a bug report via the Support Option.

However, we don't offer support for integrating LicenseChain into your project. If you’re having trouble, you might want to have a look on Google or YouTube for tutorials on the programming language you're using to build your programme.

# Copyright License
LicenseChain is under the Elastic License 2.0.

- You’re not allowed to offer the software to third parties as a hosted or managed service, where users get access to any significant part of the software’s features or functionality.
- You mustn’t move, alter, disable, or bypass the licence key functionality in the software, and you can’t remove or hide any features protected by the licence key.
- You’re also not permitted to change, remove, or obscure any licensing, copyright, or other notices from the licensor within the software. Any use of the licensor’s trademarks must comply with relevant laws.

Cheers for sticking to these guidelines. We put a lot of effort into developing LicenseChain and don't take copyright breaches lightly.

## Support

If you have any questions or need help, feel free to open an issue or contact us at support@licensechain.app.
