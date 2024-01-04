// manage base64 encoding and decoding
package cookies

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
)

var (
	ErrValueTooLong = errors.New("the current cookie value is too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)

// encodes the cookie value
func Write(w http.ResponseWriter, cookie http.Cookie) error {

	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	//check if len is too long > 4096 bytes
	if len(cookie.String()) > 4096 {
		return ErrValueTooLong
	}

	http.SetCookie(w, &cookie)
	return nil

}

// finds the cookie provided on the http req and decodes it.
func Read(r *http.Request, name string) (string, error) {

	cookie, err := r.Cookie(name)

	if err != nil {
		return "", err
	}

	cookieValue, err := base64.URLEncoding.DecodeString(cookie.Value)

	if err != nil {
		return "", ErrInvalidValue
	}

	return string(cookieValue), nil

}

func WriteSigned(w http.ResponseWriter, cookie http.Cookie, secretKey []byte) error {
	// Calculate a HMAC signature of the cookie name and value, using SHA256 and a secret key.
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(cookie.Name))
	mac.Write([]byte(cookie.Value))
	signature := mac.Sum(nil)

	// Prepend the cookie value with the HMAC signature.
	cookie.Value = string(signature) + cookie.Value

	// Call  Write() helper to base64-encode the new cookie value and write
	// the cookie.
	return Write(w, cookie)
}

func ReadSigned(r *http.Request, name string, secretKey []byte) (string, error) {
	// Read in the signed value from the cookie. This should be in the format
	// "{signature}{original value}".
	signedValue, err := Read(r, name)
	if err != nil {
		return "", err
	}

	// A SHA256 HMAC signature has a fixed length of 32 bytes. To avoid a potential
	// 'index out of range' panic in the next step, we need to check sure that the
	// length of the signed cookie value is at least this long. We'll use the
	// sha256.Size constant here, rather than 32, just because it makes our code
	// a bit more understandable at a glance.
	if len(signedValue) < sha256.Size {
		return "", ErrInvalidValue
	}

	// Split apart the signature and original cookie value.
	signature := signedValue[:sha256.Size]
	value := signedValue[sha256.Size:]

	// Recalculate the HMAC signature of the cookie name and original value.
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(name))
	mac.Write([]byte(value))
	expectedSignature := mac.Sum(nil)

	// Check that the recalculated signature matches the signature we received
	// in the cookie. If they match, we can be confident that the cookie name
	// and value haven't been edited by the client.
	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", ErrInvalidValue
	}

	// Return the original cookie value.
	return value, nil
}
