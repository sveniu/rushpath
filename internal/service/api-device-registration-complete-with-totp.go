package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type APIAuthenticationCompleteDeviceRegistrationWithTOTPInputDevice struct {
	DeviceName *string `json:"deviceName"`
	AppVersion *string `json:"appVersion"`
	Platform   *string `json:"platform"`
	OSCountry  *string `json:"osCountry"`
	OSLanguage *string `json:"osLanguage"`
	Temporary  *bool   `json:"temporary"`
}

type APIAuthenticationCompleteDeviceRegistrationWithTOTPInputVerification struct {
	OTP *string `json:"otp"`
}

type APIAuthenticationCompleteDeviceRegistrationWithTOTPInput struct {
	Login        *string                                                               `json:"login"`
	Device       *APIAuthenticationCompleteDeviceRegistrationWithTOTPInputDevice       `json:"device"`
	Verification *APIAuthenticationCompleteDeviceRegistrationWithTOTPInputVerification `json:"verification"`
}

type APIAuthenticationCompleteDeviceRegistrationWithTOTPOutputData struct {
	DeviceAccessKey *string `json:"deviceAccessKey"`
	DeviceSecretKey *string `json:"deviceSecretKey"`
}

type APIAuthenticationCompleteDeviceRegistrationWithTOTPOutput struct {
	RequestID *string                                                        `json:"requestId"`
	Data      *APIAuthenticationCompleteDeviceRegistrationWithTOTPOutputData `json:"data"`
}

func (s *Service) APIAuthenticationCompleteDeviceRegistrationWithTOTP(
	input *APIAuthenticationCompleteDeviceRegistrationWithTOTPInput,
) (
	*APIAuthenticationCompleteDeviceRegistrationWithTOTPOutput,
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
		"https://api.dashlane.com/v1/authentication/CompleteDeviceRegistrationWithTOTP",
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
	apiAuthenticationcompleteDeviceRegistrationWithTOTPOutput := &APIAuthenticationCompleteDeviceRegistrationWithTOTPOutput{}
	if err := json.NewDecoder(resp.Body).Decode(
		apiAuthenticationcompleteDeviceRegistrationWithTOTPOutput,
	); err != nil {
		return nil, err
	}

	return apiAuthenticationcompleteDeviceRegistrationWithTOTPOutput, nil
}
