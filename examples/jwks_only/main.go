// Example: verify an existing license_token with JWKS only (no POST /v1/licenses/verify).
// Use when your app already holds a token and license_jwks_uri from a prior verify response.
//
// Run from repo root:
//
//	LICENSECHAIN_LICENSE_TOKEN=eyJ... LICENSECHAIN_LICENSE_JWKS_URI=https://api.licensechain.app/v1/licenses/jwks go run ./examples/jwks_only
//
// Optional: LICENSECHAIN_EXPECTED_APP_ID=<app-uuid> to enforce JWT aud.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/licensechain/licensechain-go-sdk/client"
)

func main() {
	token := os.Getenv("LICENSECHAIN_LICENSE_TOKEN")
	jwksURL := os.Getenv("LICENSECHAIN_LICENSE_JWKS_URI")
	if token == "" || jwksURL == "" {
		log.Fatal("Set LICENSECHAIN_LICENSE_TOKEN and LICENSECHAIN_LICENSE_JWKS_URI (from verify response)")
	}

	opts := (*client.VerifyLicenseAssertionOptions)(nil)
	if app := os.Getenv("LICENSECHAIN_EXPECTED_APP_ID"); app != "" {
		opts = &client.VerifyLicenseAssertionOptions{ExpectedAppID: app}
	}

	claims, err := client.VerifyLicenseAssertionJWT(token, jwksURL, opts)
	if err != nil {
		log.Fatalf("jwt: %v", err)
	}
	fmt.Printf("token_use=%v sub=%v lc_vt=%v\n", claims["token_use"], claims["sub"], claims["lc_vt"])
}
