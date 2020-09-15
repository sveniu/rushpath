package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	signatureAlgorithm = "DL1-HMAC-SHA256"
)

func (s *Service) addAuthorizationHeader(authType authType, req *http.Request) error {
	timestamp := time.Now().Unix()

	// Build canonical request components.
	canonicalURI := req.URL.Path

	canonicalQueryString := "" // FIXME implement this

	canonicalHeaders := strings.Join([]string{
		"accept:" + req.Header.Get("accept"),
		"content-type:" + req.Header.Get("content-type"),
		"dashlane-client-agent:" + req.Header.Get("dashlane-client-agent"),
	}, "\n") + "\n"

	signedHeaders := strings.Join([]string{
		"accept",
		"content-type",
		"dashlane-client-agent",
	}, ";")

	bodyReadCloser, err := req.GetBody()
	if err != nil {
		return err
	}

	bodyBytes, err := ioutil.ReadAll(bodyReadCloser)
	if err != nil {
		return err
	}

	payloadHash := hex.EncodeToString(hashSHA256(bodyBytes))

	canonicalRequest := strings.Join([]string{
		req.Method,
		canonicalURI,
		canonicalQueryString,
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	}, "\n")

	stringToSign := strings.Join([]string{
		signatureAlgorithm,
		fmt.Sprintf("%d", timestamp),
		hex.EncodeToString(hashSHA256([]byte(canonicalRequest))),
	}, "\n")

	signingKey := *s.Credentials.AppSecretKey
	switch authType {
	case authTypeUserDevice:
		signingKey += "\n" + *s.Credentials.DeviceSecretKey
	}

	signature := hex.EncodeToString(hmacSHA256([]byte(signingKey), []byte(stringToSign)))

	headerValue := signatureAlgorithm + " "
	switch authType {
	case authTypeApp:
		headerValue += "AppAccessKey=" + *s.Credentials.AppAccessKey
	case authTypeUserDevice:
		headerValue += "Login=" + *s.Credentials.Login
		headerValue += ",AppAccessKey=" + *s.Credentials.AppAccessKey
		headerValue += ",DeviceAccessKey=" + *s.Credentials.DeviceAccessKey
	}
	headerValue += ",Timestamp=" + fmt.Sprintf("%d", timestamp)
	headerValue += ",SignedHeaders=" + signedHeaders
	headerValue += ",Signature=" + signature

	req.Header.Set("authorization", headerValue)

	log.Trace().
		Bytes("payload", bodyBytes).
		Str("payload_hash", payloadHash).
		Str("canonical_request", canonicalRequest).
		Str("string_to_sign", stringToSign).
		Str("authorization_header", headerValue).
		Msg("signing request")

	return nil
}

func hmacSHA256(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func hashSHA256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}
