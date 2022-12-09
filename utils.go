package cc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

type tokenHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func camel2Case(s string) string {
	var output []rune
	for i, r := range s {
		if i == 0 {
			output = append(output, unicode.ToLower(r))
		} else {
			if unicode.IsUpper(r) {
				output = append(output, '_')
			}
			output = append(output, unicode.ToLower(r))
		}
	}
	return string(output)
}

func JWTToken(payload any, secret string) (string, error) {
	header := tokenHeader{Alg: "HS256", Typ: "JWT"}
	json_header, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	base64_header := base64.StdEncoding.EncodeToString(json_header)
	json_payload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	base64_payload := base64.StdEncoding.EncodeToString(json_payload)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(fmt.Sprintf("%s.%s", base64_header, base64_payload)))
	signature := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s.%s.%s", base64_header, base64_payload, signature), nil
}

func JWTCheck(token string, secret string) (bool, error) {
	infos := strings.Split(token, ".")
	sig, err := hex.DecodeString(infos[2])
	if err != nil {
		return false, err
	}
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(fmt.Sprintf("%s.%s", infos[0], infos[1])))
	if !hmac.Equal(sig, h.Sum(nil)) {
		return false, nil
	}
	return true, nil
}

func JWTParse(token string, result any) error {
	infos := strings.Split(token, ".")
	payload, err := base64.StdEncoding.DecodeString(infos[1])
	if err != nil {
		return err
	}
	err = json.Unmarshal(payload, result)
	if err != nil {
		return err
	}
	return nil
}
