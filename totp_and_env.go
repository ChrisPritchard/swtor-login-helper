package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

func generateTOTP(secret string) (int, error) {
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return 0, fmt.Errorf("invalid base32 secret: %v", err)
	}

	counter := time.Now().Unix() / 30

	h := hmac.New(sha1.New, key)
	binary.Write(h, binary.BigEndian, counter)
	hash := h.Sum(nil)

	offset := hash[len(hash)-1] & 0xF
	truncatedHash := hash[offset : offset+4]

	code := binary.BigEndian.Uint32(truncatedHash) & 0x7FFFFFFF
	code = code % uint32(math.Pow10(6))

	return int(code), nil
}

func readEnvFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	envVars := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // skip malformed lines
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if len(value) > 0 && (value[0] == '"' || value[0] == '\'') {
			value = value[1 : len(value)-1]
		}

		envVars[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envVars, nil
}
