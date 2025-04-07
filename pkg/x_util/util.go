package x_util

import (
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/rskv-p/jtask/pkg/x_log"
)

// charset defines allowed characters for the random string.
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString generates a secure random alphanumeric string of the given length.
// Returns an error if secure randomness could not be obtained.
func RandomString(length int) (string, error) {
	// Preallocate a buffer for bytes and the result string builder.
	buf := make([]byte, length)
	result := strings.Builder{}
	result.Grow(length)

	// Fill the buffer with cryptographically secure random bytes.
	if _, err := rand.Read(buf); err != nil {
		// Log error if randomness generation fails
		x_log.Error().
			Err(err).
			Int("length", length).
			Msg("failed to generate secure random string")
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}

	// Map each byte to a character in the charset.
	for i := 0; i < length; i++ {
		result.WriteByte(charset[int(buf[i])%len(charset)])
	}

	// Log success for random string generation
	x_log.Debug().
		Int("length", length).
		Str("random_string", result.String()).
		Msg("random string generated successfully")

	return result.String(), nil
}
