// Example: POST /v1/licenses/verify then verify license_token with JWKS (RS256).
//
// Run from repo root:
//
//	LICENSECHAIN_API_KEY=... LICENSECHAIN_LICENSE_KEY=... go run ./examples/license_assertion
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/licensechain/licensechain-go-sdk/client"
)

func main() {
	apiKey := os.Getenv("LICENSECHAIN_API_KEY")
	licenseKey := os.Getenv("LICENSECHAIN_LICENSE_KEY")
	base := os.Getenv("LICENSECHAIN_BASE_URL")
	if base == "" {
		base = "https://api.licensechain.app/v1"
	}
	if apiKey == "" || licenseKey == "" {
		log.Fatal("Set LICENSECHAIN_API_KEY and LICENSECHAIN_LICENSE_KEY (optional LICENSECHAIN_BASE_URL)")
	}

	lc := client.CreateClient(apiKey, base)

	details, err := lc.VerifyLicenseWithDetails(licenseKey)
	if err != nil {
		log.Fatalf("verify: %v", err)
	}
	valid, _ := details["valid"].(bool)
	if !valid {
		log.Fatalf("license not valid: %#v", details)
	}
	token, _ := details["license_token"].(string)
	jwksURL, _ := details["license_jwks_uri"].(string)
	if token == "" || jwksURL == "" {
		log.Print("No license_token or license_jwks_uri — enable LICENSE_JWT_* on Core API for this seller.")
		return
	}

	claims, err := client.VerifyLicenseAssertionJWT(token, jwksURL, nil)
	if err != nil {
		log.Fatalf("jwt: %v", err)
	}
	fmt.Printf("token_use=%v sub=%v lc_vt=%v\n", claims["token_use"], claims["sub"], claims["lc_vt"])
}
