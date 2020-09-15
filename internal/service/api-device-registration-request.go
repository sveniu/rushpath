package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type APIRequestDeviceRegistrationInput struct {
	Login *string `json:"login"`
}

type APIRequestDeviceRegistrationOutputDataVerificationU2FChallenge struct {
	Version   *string `json:"version"`
	KeyHandle *string `json:"keyHandle"`
	Challenge *string `json:"challenge"`
	AppID     *string `json:"appId"`
}

type APIRequestDeviceRegistrationOutputDataVerification struct {
	Type       *string                                                           `json:"type"`
	Challenges []*APIRequestDeviceRegistrationOutputDataVerificationU2FChallenge `json:"challenges,omitempty"`
}

type APIRequestDeviceRegistrationOutputData struct {
	Verification []*APIRequestDeviceRegistrationOutputDataVerification `json:"verification"`
}

type APIRequestDeviceRegistrationOutput struct {
	RequestID *string                                 `json:"requestId"`
	Data      *APIRequestDeviceRegistrationOutputData `json:"data"`
}

func (s *Service) APIRequestDeviceRegistration(
	input *APIRequestDeviceRegistrationInput,
) (
	*APIRequestDeviceRegistrationOutput,
	error,
) {
	// Create client if it doesn't exist already.
	if s.Client == nil {
		s.Client = &http.Client{}
	}

	// Prepare request body.
	bodyBuf := new(bytes.Buffer)
	compactBuf := new(bytes.Buffer)
	json.NewEncoder(bodyBuf).Encode(input)
	if err := json.Compact(compactBuf, bodyBuf.Bytes()); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.dashlane.com/v1/authentication/RequestDeviceRegistration",
		compactBuf,
	)
	if err != nil {
		return nil, err
	}

	// Set headers.
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("dashlane-client-agent", `{"platform":"server_standalone"}`)

	// Sign the request.
	s.addAuthorizationHeader(authTypeApp, req)

	// Execute the request.
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Not-OK return code: %d", resp.StatusCode)
	}

	// Decode the response.
	apiRequestDeviceRegistrationOutput := &APIRequestDeviceRegistrationOutput{}
	if err := json.NewDecoder(resp.Body).Decode(
		apiRequestDeviceRegistrationOutput,
	); err != nil {
		return nil, err
	}

	return apiRequestDeviceRegistrationOutput, nil
}
