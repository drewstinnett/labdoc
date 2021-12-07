package labdoc

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// Get sha256 from file.
func GetSha256(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	hs := fmt.Sprintf("%x", h.Sum(nil))
	return hs, nil
}
