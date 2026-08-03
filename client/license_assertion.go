package client

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

// LICENSE_TOKEN_USE_CLAIM must match Core API license JWTs.
const LICENSE_TOKEN_USE_CLAIM = "licensechain_license_v1"

// VerifyLicenseAssertionOptions optional checks after RS256 + JWKS verification.
type VerifyLicenseAssertionOptions struct {
	ExpectedAppID string // when non-empty, JWT aud must match
	Issuer        string // when non-empty, passed to jwt parser as issuer
}

// VerifyLicenseAssertionJWT verifies license_token from POST /v1/licenses/verify using JWKS (RS256).
func VerifyLicenseAssertionJWT(token string, jwksURL string, opts *VerifyLicenseAssertionOptions) (jwt.MapClaims, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errors.New("empty token")
	}
	if strings.TrimSpace(jwksURL) == "" {
		return nil, errors.New("empty jwksURL")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	jwks, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("jwks: %w", err)
	}

	parseOpts := []jwt.ParserOption{jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()})}
	if opts != nil && opts.Issuer != "" {
		parseOpts = append(parseOpts, jwt.WithIssuer(opts.Issuer))
	}

	parsed, err := jwt.Parse(token, jwks.Keyfunc, parseOpts...)
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("unexpected claims type")
	}
	tu, _ := claims["token_use"].(string)
	if tu != LICENSE_TOKEN_USE_CLAIM {
		return nil, fmt.Errorf("token_use: want %q", LICENSE_TOKEN_USE_CLAIM)
	}
	if opts != nil && strings.TrimSpace(opts.ExpectedAppID) != "" {
		aud := claims["aud"]
		switch v := aud.(type) {
		case string:
			if v != opts.ExpectedAppID {
				return nil, errors.New("aud mismatch")
			}
		case []interface{}:
			ok := false
			for _, x := range v {
				if s, ok2 := x.(string); ok2 && s == opts.ExpectedAppID {
					ok = true
					break
				}
			}
			if !ok {
				return nil, errors.New("aud mismatch")
			}
		default:
			if aud != opts.ExpectedAppID {
				return nil, errors.New("aud mismatch")
			}
		}
	}
	return claims, nil
}
