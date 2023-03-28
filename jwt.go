package cc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// JwtToken JWT加密获取token
func JWTToken(payload map[string]any, secret []byte) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	header_str, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	header_result := base64.URLEncoding.EncodeToString(header_str)
	payload_str, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	payload_result := base64.URLEncoding.EncodeToString(payload_str)
	hs := hmac.New(sha256.New, secret)
	if _, err := hs.Write([]byte(fmt.Sprintf("%s.%s", header_result, payload_result))); err != nil {
		return "", err
	}
	sum := base64.URLEncoding.EncodeToString(hs.Sum(nil))
	return fmt.Sprintf("%s.%s.%s", header_result, payload_result, sum), nil
}

// ParseJWT 解析JWT
func ParseJWT(token string, secret []byte) (map[string]any, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("token format error")
	}
	hs := hmac.New(sha256.New, secret)
	if _, err := hs.Write([]byte(fmt.Sprintf("%s.%s", parts[0], parts[1]))); err != nil {
		return nil, fmt.Errorf("encrypt error")
	}
	if base64.URLEncoding.EncodeToString(hs.Sum(nil)) != parts[2] {
		return nil, fmt.Errorf("token is invalid")
	}
	var result map[string]any
	result_bytes, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("token decode error")
	}
	if err := json.Unmarshal(result_bytes, &result); err != nil {
		return nil, fmt.Errorf("token format error")
	}
	return result, nil
}
