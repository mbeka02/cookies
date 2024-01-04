// manage base64 encoding and decoding
package cookies

import (
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
//finds the cookie provided on the http req and decodes it.
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
