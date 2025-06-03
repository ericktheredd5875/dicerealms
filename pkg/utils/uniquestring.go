package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"time"
)

func GenerateUniqueString(byteLength int) (string, error) {
	b := make([]byte, byteLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return hex.EncodeToString(b), nil
}

func GenerateUniqueID() (string, error) {
	return GenerateUniqueString(16)
}

func GenerateKHash(input string, dataType string) (string, error) {
	if input == "" {
		var err error
		input, err = GenerateUniqueID()
		if err != nil {
			return "", fmt.Errorf("failed to generate KHash: %w", err)
		}
	}

	var timeInput string
	if dataType == "date" {
		timeInput = time.Now().Format("2006-01-02")
	} else {
		timeInput = fmt.Sprintf("%d", time.Now().Unix())
	}

	fullInput := input + ":" + timeInput
	crc := crc32.ChecksumIEEE([]byte(fullInput))
	return fmt.Sprintf("%08x", crc), nil
}
