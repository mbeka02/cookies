package main

import (
	"encoding/hex"
	"errors"
	"log"
	"net/http"

	"github.com/mbeka02/cookies_go/internal/cookies"
)

var secretKey []byte

func main() {
	var err error
	// Decode the random 64-character hex string to give us a slice containing 32 random bytes
	secretKey, err = hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		log.Fatal(err)
	}
	var port = "3000"
	mux := http.NewServeMux()
	mux.HandleFunc("/set", setCookieHandler)
	mux.HandleFunc("/get", getCookieHandler)
	log.Printf("Listening on port : %v", port)
	err = http.ListenAndServe(":"+port, mux)

	if err != nil {
		log.Fatal(err)
	}
}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {

	cookie := http.Cookie{
		Name:     "DemoCookie",
		Value:    "12345686790",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	err := cookies.WriteSigned(w, cookie, secretKey)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

	}

	w.Write([]byte("cookie has been set"))

}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {

	val, err := cookies.ReadSigned(r, "DemoCookie", secretKey)

	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusBadRequest)
		case errors.Is(err, cookies.ErrInvalidValue):
			http.Error(w, "invalid cookie", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}

		w.Write([]byte(val))

	}

}
