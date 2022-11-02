// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package sessions

import (
	"crypto/sha256"
	"encoding/base64"
)

// Creates CSRF token for session.
// Input should be base64 encoded.
func CSRFToken(sessionID string) string {
	bytes, err := base64.StdEncoding.DecodeString(sessionID)
	if err != nil {
		panic(err)
	}
	result := sha256.Sum256(bytes)
	bytes = result[:]
	return base64.StdEncoding.EncodeToString(bytes)
}

// Validates CSRF token.
func CheckCSRFToken(sessionID, token string) bool {
	return CSRFToken(sessionID) == token
}
