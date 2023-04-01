package pow

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

const (
	zeroByte   rune   = 48             // ASCII code for number zero
	timeFormat string = "060102150405" // YYMMDDhhmmss
)

// HashcashData - struct with fields of Hashcash
type HashcashData struct {
	// version hashcash format version, 1 (which supersedes version 0).
	Version int
	// bits number of "partial pre-image" (zero) bits in the hashed code.
	Bits int
	// created date The time that the message was sent.
	Date time.Time
	// resource data string being transmitted, e.g., an IP address or email address.
	Resource string
	// extension (optional; ignored in version 1).
	Extension string
	// rand characters, encoded in base-64 format.
	Rand string
	// counter (up to 2^20), encoded in base-64 format.
	Counter int
}

// createHeader - creates a new hashcash header
func (h HashcashData) createHeader() string {
	return fmt.Sprintf("%d:%d:%s:%s:%s:%s:%d", h.Version, h.Bits, h.Date.Format(timeFormat), h.Resource, h.Extension, h.Rand, h.Counter)
}

// sha256Hash - calculates sha256 hash from given string
func sha256Hash(data string) string {
	h := sha256.New()
	h.Write([]byte(data))

	return fmt.Sprintf("%x", h.Sum(nil))
}

// Valid - checks that hash is valid
func IsHashValid(hash string, zerosCount int) bool {
	if len(hash) < zerosCount {
		return false
	}
	for _, ch := range hash[:zerosCount] {
		if ch != zeroByte {
			return false
		}
	}
	return true
}

// CalculateHashcash - calculates correct hashcash by bruteforce
func (h HashcashData) CalculateHashcash(maxIterations int) (HashcashData, error) {
	for h.Counter <= maxIterations || maxIterations <= 0 {
		header := h.createHeader()
		hash := sha256Hash(header)
		if IsHashValid(hash, h.Bits) {
			return h, nil
		}
		// if hash don't have needed count of leading zeros, we are increasing counter and try next hash
		h.Counter++
	}

	return h, fmt.Errorf("max iterations reached")
}

// NewHashcash creates a new Hashcash
func NewHashcash(resource string, rand int, zeroCnt int) (*HashcashData, error) {
	if resource == "" {
		return nil, fmt.Errorf("resource is empty")
	}

	return &HashcashData{
		Version:  1,
		Bits:     zeroCnt,
		Date:     time.Now().UTC().UTC(),
		Resource: resource,
		Rand:     base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", rand))),
		Counter:  0,
	}, nil
}

// CalculateeHashcash - calculates correct hashcash by bruteforce
func (h HashcashData) GetRand() (int, error) {
	data, err := base64.StdEncoding.DecodeString(h.Rand)
	if err != nil {
		return 0, err
	}

	randValue, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}

	return randValue, nil
}
