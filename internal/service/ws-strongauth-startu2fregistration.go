package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

type WSStrongauthStartU2FRegistrationInput struct {
	Login *string `url:"login"`
	UKI   *string `url:"uki"`
}

type WSStrongauthStartU2FRegistrationOutputChallenge struct {
	Version   *string `json:"version"`
	AppID     *string `json:"appId"`
	Challenge *string `json:"challenge"`
}

type WSStrongauthStartU2FRegistrationOutput struct {
	Code    *int                                             `json:"code"`
	Message *string                                          `json:"message"`
	Content *WSStrongauthStartU2FRegistrationOutputChallenge `json:"content"`
}

func (s *Service) WSStrongauthStartU2FRegistration(
	input *WSStrongauthStartU2FRegistrationInput,
) (
	*WSStrongauthStartU2FRegistrationOutput,
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
		"https://ws1.dashlane.com/3/strongauth/startU2FRegistration",
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
	wsStrongauthStartU2FRegistrationOutput := &WSStrongauthStartU2FRegistrationOutput{}
	if err := json.NewDecoder(resp.Body).Decode(
		wsStrongauthStartU2FRegistrationOutput,
	); err != nil {
		return nil, err
	}

	return wsStrongauthStartU2FRegistrationOutput, nil
}
