package adyen_encrypt

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	chunkSize = 32768
	re        = regexp.MustCompile(`(\d{4})(\d{4})(\d{4})(\d{4})`)
)

func NowTimeISO() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.999") + "Z"
}

func EncodeToBase64(input interface{}) string {
	var inputData []byte

	switch input := input.(type) {
	case string:
		inputData = []byte(input)
	case []byte:
		inputData = input
	default:
		fmt.Println("Adyen Encrypt: [EncodeToBase64: unsupported input type]")
		return ""
	}

	var chunks []string
	for n := 0; n < len(inputData); n += chunkSize {
		end := n + chunkSize
		if end > len(inputData) {
			end = len(inputData)
		}
		chunks = append(chunks, string(inputData[n:end]))
	}

	combinedData := []byte{}
	for _, chunk := range chunks {
		combinedData = append(combinedData, chunk...)
	}

	encoded := base64.StdEncoding.EncodeToString(combinedData)
	encoded = strings.ReplaceAll(encoded, "=", "")
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")
	return encoded
}

func HexDecode(e string) []byte {
	if e == "" {
		return []byte{}
	}

	if len(e)%2 == 1 {
		e = "0" + e
	}

	t := len(e) / 2
	r := make([]byte, t)

	for n := 0; n < t; n++ {
		hexByte := e[2*n : 2*n+2]
		byteValue, err := hex.DecodeString(hexByte)
		if err != nil {
			fmt.Println("Error decoding hex string:", err)
			return nil
		}
		r[n] = byteValue[0]
	}

	return r
}

func FormatCardNumber(cardNumber string) string {
	formatted := re.ReplaceAllString(cardNumber, "$1 $2 $3 $4")
	return formatted
}
