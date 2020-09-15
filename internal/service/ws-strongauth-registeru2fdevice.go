package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

type WSStrongauthRegisterU2FDeviceInput struct {
	ChallengeAnswer *string `url:"challengeAnswer"`
	Name            *string `url:"name"`
	Login           *string `url:"login"`
	UKI             *string `url:"uki"`
}

type WSStrongauthRegisterU2FDeviceOutput struct {
	Code    *int    `json:"code"`
	Message *string `json:"message"`
}

func (s *Service) WSStrongauthRegisterU2FDevice(
	input *WSStrongauthRegisterU2FDeviceInput,
) (
	*WSStrongauthRegisterU2FDeviceOutput,
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
		"https://ws1.dashlane.com/3/strongauth/registerU2FDevice",
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
	wsStrongauthRegisterU2FDeviceOutput := &WSStrongauthRegisterU2FDeviceOutput{}
	if err := json.NewDecoder(resp.Body).Decode(
		wsStrongauthRegisterU2FDeviceOutput,
	); err != nil {
		return nil, err
	}

	return wsStrongauthRegisterU2FDeviceOutput, nil
}
