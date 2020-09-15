package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

type WSStrongauthUnregisterU2FDeviceInput struct {
	KeyHandle *string `url:"keyHandle"`
	Login     *string `url:"login"`
	UKI       *string `url:"uki"`
}

type WSStrongauthUnregisterU2FDeviceOutput struct {
	Code    *int    `json:"code"`
	Message *string `json:"message"`
}

func (s *Service) WSStrongauthUnregisterU2FDevice(
	input *WSStrongauthUnregisterU2FDeviceInput,
) (
	*WSStrongauthUnregisterU2FDeviceOutput,
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
		"https://ws1.dashlane.com/3/strongauth/unregisterU2FDevice",
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
	wsStrongauthUnregisterU2FDeviceOutput := &WSStrongauthUnregisterU2FDeviceOutput{}
	if err := json.NewDecoder(resp.Body).Decode(
		wsStrongauthUnregisterU2FDeviceOutput,
	); err != nil {
		return nil, err
	}

	return wsStrongauthUnregisterU2FDeviceOutput, nil
}
