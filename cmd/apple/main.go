package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type AppleIdentityClaims struct {
	*jwt.StandardClaims
	Email          string `json:"email"`
	EmailVerified  string `json:"email_verified"`
	IsPrivateEmail string `json:"is_private_email"`
	SubjectID      string `json:"sub"`
}

func ValidateAppleIdentityToken(
	identityToken string,
) (*AppleIdentityClaims, error) {
	// Get Apple's public key
	resp, err := http.Get("https://appleid.apple.com/auth/keys")
	if err != nil {
		return nil, fmt.Errorf("error getting Apple's public key: %v", err)
	}
	defer resp.Body.Close()

	// Parse the token
	token, err := jwt.ParseWithClaims(
		identityToken,
		&AppleIdentityClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Verify signing algorithm
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Header["alg"],
				)
			}

			// Get the key ID from token header
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, fmt.Errorf("kid header not found")
			}

			// Read Apple's public keys
			keyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf(
					"error reading Apple's public key response: %v",
					err,
				)
			}

			var keys struct {
				Keys []struct {
					Kty string `json:"kty"`
					Kid string `json:"kid"`
					Use string `json:"use"`
					Alg string `json:"alg"`
					N   string `json:"n"`
					E   string `json:"e"`
				} `json:"keys"`
			}

			if err := json.Unmarshal(keyBytes, &keys); err != nil {
				return nil, fmt.Errorf(
					"error parsing Apple's public keys: %v",
					err,
				)
			}

			// Find the correct key
			var publicKey interface{}
			for _, key := range keys.Keys {
				if key.Kid == kid {
					publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(key.N))
					if err != nil {
						return nil, fmt.Errorf(
							"error parsing RSA public key: %v",
							err,
						)
					}
					break
				}
			}

			if publicKey == nil {
				return nil, fmt.Errorf("matching key not found")
			}

			return publicKey, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("error parsing identity token: %v", err)
	}

	if claims, ok := token.Claims.(*AppleIdentityClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func main() {
	identityToken := "YOUR_APPLE_IDENTITY_TOKEN"

	claims, err := ValidateAppleIdentityToken(identityToken)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("User ID: %s\n", claims.SubjectID)
	fmt.Printf("Email: %s\n", claims.Email)
	fmt.Printf("Email Verified: %s\n", claims.EmailVerified)
	fmt.Printf("Is Private Email: %s\n", claims.IsPrivateEmail)
}
