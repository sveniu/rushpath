package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

type WSStrongauthDeactivateForNewDeviceInput struct {
	Login *string `url:"login"`
	UKI   *string `url:"uki"`
	OTP   *string `url:"otp"`
}

type WSStrongauthDeactivateForNewDeviceOutput struct {
	Code    *int    `json:"code"`
	Message *string `json:"message"`
}

func (s *Service) WSStrongauthDeactivateForNewDevice(
	input *WSStrongauthDeactivateForNewDeviceInput,
) (
	*WSStrongauthDeactivateForNewDeviceOutput,
	error,
) {
	// Create client if it doesn't exist already.
	if s.Client == nil {
		s.Client = &http.Client{}
	}

	// Set up form parameters.
	formParams, err := query.Values(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		"https://ws1.dashlane.com/3/strongauth/deactivateForNewDevice",
		strings.NewReader(formParams.Encode()),
	)
	if err != nil {
		return nil, err
	}

	// Execute the request.
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Not-OK return code: %d", resp.StatusCode)
	}

	// Decode the response.
	wsStrongauthDeactivateForNewDeviceOutput := &WSStrongauthDeactivateForNewDeviceOutput{}
	if err := json.NewDecoder(resp.Body).Decode(
		wsStrongauthDeactivateForNewDeviceOutput,
	); err != nil {
		return nil, err
	}

	return wsStrongauthDeactivateForNewDeviceOutput, nil
}
