package app

import (
	"crypto/md5"
	"fmt"
)

// createSignature will return a signature for the given response set.
//
// Note: probably it should use RSA to sign this, but I'm way over time already.
func createSignature(rs ResponseSet, secret string) (string, error) {
	h := md5.New()
	h.Write([]byte(rs.Email))
	h.Write([]byte{0})
	for _, r := range rs.Responses {
		h.Write([]byte(r.Question))
		h.Write([]byte{0})
		h.Write([]byte(r.Answer))
		h.Write([]byte{0})
	}
	h.Write([]byte(secret))
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
