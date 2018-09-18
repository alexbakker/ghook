// Package ghook provides a minimal toolset for receiving with GitHub web hooks.
package ghook

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type (
	// Callback is the type of function that is called when a GitHub web hook
	// request is received and verified.
	Callback func(event *Event) error

	Hook struct {
		secret []byte
		cb     Callback
	}

	Event struct {
		Name    string
		GUID    string
		Payload []byte
	}
)

func New(secret []byte, cb Callback) *Hook {
	return &Hook{secret: secret, cb: cb}
}

// ServeHTTP implements the http.Handler interface.
func (h *Hook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		writeError(w, errors.New("bad content type"), http.StatusBadRequest)
		return
	}

	name, guid, digest, err := parseHeader(r)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	digestBytes, err := hex.DecodeString(digest)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	mac := hmac.New(sha1.New, h.secret)
	tr := io.TeeReader(r.Body, mac)

	data, err := ioutil.ReadAll(tr)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	if !hmac.Equal(mac.Sum(nil), digestBytes) {
		writeError(w, errors.New("bad digest"), http.StatusForbidden)
		return
	}

	event := Event{
		Name:    name,
		GUID:    guid,
		Payload: data,
	}

	if err = h.cb(&event); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
}

// parseHeader parses the GitHub web hook header of the given request and
// returns the digest and event type.
func parseHeader(r *http.Request) (event string, guid string, digest string, err error) {
	if event, err = readHeader(r, "X-GitHub-Event"); err != nil {
		return
	}
	if guid, err = readHeader(r, "X-GitHub-Delivery"); err != nil {
		return
	}
	if digest, err = readHeader(r, "X-Hub-Signature"); err != nil {
		return
	}

	return event, guid, strings.TrimLeft(digest, "sha1="), nil
}

// readHeader retrieves the value of the given key from the header.
func readHeader(r *http.Request, key string) (string, error) {
	value := r.Header.Get(key)
	if value == "" {
		return "", fmt.Errorf("%s is not set", key)
	}

	return value, nil
}

// writeStatus writes the given HTTP status code to the given
// http.ResponseWriter.
func writeError(w http.ResponseWriter, err error, status int) {
	http.Error(w, fmt.Sprintf("error: %s", err), status)
}
