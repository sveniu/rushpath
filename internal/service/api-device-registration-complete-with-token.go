package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type APIAuthenticationCompleteDeviceRegistrationWithTokenInputDevice struct {
	DeviceName *string `json:"deviceName"`
	AppVersion *string `json:"appVersion"`
	Platform   *string `json:"platform"`
	OSCountry  *string `json:"osCountry"`
	OSLanguage *string `json:"osLanguage"`
	Temporary  *bool   `json:"temporary"`
}

type APIAuthenticationCompleteDeviceRegistrationWithTokenInputVerification struct {
	Token *string `json:"token"`
}

type APIAuthenticationCompleteDeviceRegistrationWithTokenInput struct {
	Login        *string                                                                `json:"login"`
	Device       *APIAuthenticationCompleteDeviceRegistrationWithTokenInputDevice       `json:"device"`
	Verification *APIAuthenticationCompleteDeviceRegistrationWithTokenInputVerification `json:"verification"`
}

type APIAuthenticationCompleteDeviceRegistrationWithTokenOutputData struct {
	DeviceAccessKey *string `json:"deviceAccessKey"`
	DeviceSecretKey *string `json:"deviceSecretKey"`
}

type APIAuthenticationCompleteDeviceRegistrationWithTokenOutput struct {
	RequestID *string                                                         `json:"requestId"`
	Data      *APIAuthenticationCompleteDeviceRegistrationWithTokenOutputData `json:"data"`
}

func (s *Service) APIAuthenticationCompleteDeviceRegistrationWithToken(
	input *APIAuthenticationCompleteDeviceRegistrationWithTokenInput,
) (
	*APIAuthenticationCompleteDeviceRegistrationWithTokenOutput,
	error,
) {
	// Create client if it doesn't exist already.
	if s.Client == nil {
		s.Client = &http.Client{}
	}

	// Prepare request body.
	bodyBuf := new(bytes.Buffer)
	json.NewEncoder(bodyBuf).Encode(input)

	req, err := http.NewRequest(
		"POST",
		"https://api.dashlane.com/v1/authentication/CompleteDeviceRegistrationWithToken",
		bodyBuf,
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
	apiAuthenticationcompleteDeviceRegistrationWithTokenOutput := &APIAuthenticationCompleteDeviceRegistrationWithTokenOutput{}
	if err := json.NewDecoder(resp.Body).Decode(
		apiAuthenticationcompleteDeviceRegistrationWithTokenOutput,
	); err != nil {
		return nil, err
	}

	return apiAuthenticationcompleteDeviceRegistrationWithTokenOutput, nil
}
