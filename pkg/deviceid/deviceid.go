// Package deviceid implements Astarte's 128-bit device identifier and its
// canonical encodings: the 22-character unpadded base64url wire form used in
// MQTT topics, REST paths and certificate CNs, and the canonical UUID form
// used by operators and tooling.
//
// Astarte device IDs are opaque 128-bit values. In practice they are either
// random (UUIDv4-shaped, `astartectl utils device-id generate-random` parity)
// or derived deterministically from a namespace UUID plus an arbitrary payload
// string via UUIDv5 (RFC 4122 §4.3, `astartectl utils device-id
// compute-from-string` parity).
package deviceid

import (
	"crypto/rand"
	// UUIDv5 (RFC 4122 §4.3) mandates SHA-1 for deterministic name-based
	// IDs; it is an identifier derivation, not a security control.
	"crypto/sha1" // #nosec G505
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

// EncodedLen is the length of the base64url wire form of a device ID:
// 16 bytes encoded as unpadded base64url.
const EncodedLen = 22

// ErrInvalid is wrapped by every parse failure returned from this package,
// so callers can classify rejection with errors.Is regardless of the
// specific cause (length, alphabet, canonicality, UUID syntax).
var ErrInvalid = errors.New("invalid device ID")

// encoding is unpadded base64url in strict mode: the trailing (unused)
// padding bits of the 22nd character must be zero, so every accepted string
// round-trips byte-identically. Upstream Astarte decodes device IDs the same
// way (Elixir Base.url_decode64!(padding: false)).
var encoding = base64.RawURLEncoding.Strict()

// ID is a 128-bit Astarte device identifier.
type ID [16]byte

// Parse decodes the Astarte wire form of a device ID: exactly 22 characters
// of unpadded base64url decoding to 16 bytes. Anything else — wrong length,
// padding, standard-alphabet characters ('+', '/'), non-canonical trailing
// bits — is rejected with an error wrapping ErrInvalid.
func Parse(s string) (ID, error) {
	var id ID
	if len(s) != EncodedLen {
		return id, fmt.Errorf("%w: got %d characters, want %d", ErrInvalid, len(s), EncodedLen)
	}
	n, err := encoding.Decode(id[:], []byte(s))
	if err != nil {
		return ID{}, fmt.Errorf("%w: not unpadded base64url: %v", ErrInvalid, err)
	}
	if n != len(id) {
		// Unreachable with a 22-character input, kept as a guard.
		return ID{}, fmt.Errorf("%w: decoded to %d bytes, want %d", ErrInvalid, n, len(id))
	}
	return id, nil
}

// String returns the 22-character unpadded base64url wire form.
func (id ID) String() string {
	return encoding.EncodeToString(id[:])
}

// FromUUID parses a canonical UUID string (8-4-4-4-12 hexadecimal groups,
// case-insensitive) into a device ID. Non-canonical forms (braces, URN
// prefix, missing hyphens) are rejected.
func FromUUID(s string) (ID, error) {
	var id ID
	if len(s) != 36 || s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return id, fmt.Errorf("%w: not a canonical UUID: %q", ErrInvalid, s)
	}
	dst := id[:]
	for _, group := range []string{s[0:8], s[9:13], s[14:18], s[19:23], s[24:36]} {
		n, err := hex.Decode(dst, []byte(group))
		if err != nil {
			return ID{}, fmt.Errorf("%w: not a canonical UUID: %q", ErrInvalid, s)
		}
		dst = dst[n:]
	}
	return id, nil
}

// UUID returns the canonical lowercase UUID form of the device ID.
func (id ID) UUID() string {
	var buf [36]byte
	hex.Encode(buf[0:8], id[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], id[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], id[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], id[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:36], id[10:16])
	return string(buf[:])
}

// Random returns a new random device ID. The ID is UUIDv4-shaped (version
// and RFC 4122 variant bits set), matching what `astartectl utils device-id
// generate-random` and the official SDKs produce.
func Random() (ID, error) {
	var id ID
	if _, err := rand.Read(id[:]); err != nil {
		return ID{}, fmt.Errorf("deviceid: gathering randomness: %w", err)
	}
	id[6] = (id[6] & 0x0f) | 0x40 // version 4
	id[8] = (id[8] & 0x3f) | 0x80 // RFC 4122 variant
	return id, nil
}

// FromNamespace derives a deterministic device ID from a namespace UUID and
// an arbitrary payload string using UUIDv5 (RFC 4122 §4.3: SHA-1 over
// namespace bytes followed by the payload, version and variant bits forced).
// It produces exactly the IDs `astartectl utils device-id
// compute-from-string <namespace-uuid> <payload>` prints.
func FromNamespace(ns ID, payload string) ID {
	// #nosec G401 -- SHA-1 is the digest RFC 4122 specifies for UUIDv5;
	// collision resistance is not relied upon for security here.
	h := sha1.New()
	h.Write(ns[:])
	h.Write([]byte(payload))
	sum := h.Sum(nil)

	var id ID
	copy(id[:], sum[:16])
	id[6] = (id[6] & 0x0f) | 0x50 // version 5
	id[8] = (id[8] & 0x3f) | 0x80 // RFC 4122 variant
	return id
}
