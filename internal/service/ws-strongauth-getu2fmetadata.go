package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

type WSStrongauthGetU2FMetadataInput struct {
	Login *string `url:"login"`
	UKI   *string `url:"uki"`
}

type WSStrongauthGetU2FMetadataOutputKey struct {
	KeyHandle           *string `json:"keyHandle"`
	Name                *string `json:"name"`
	CreationDateUnix    *int64  `json:"creationDateUnix"`
	LastUsedDateUnix    *int64  `json:"lastUsedDateUnix"`
	LastUsedFromIP      *string `json:"lastUsedFromIP"`
	LastUsedFromThisIP  *bool   `json:"lastUsedFromThisIp"`
	LastUsedFromCountry *string `json:"lastUsedFromCountry"`
}

type WSStrongauthGetU2FMetadataOutput struct {
	Code    *int                                   `json:"code"`
	Message *string                                `json:"message"`
	Content []*WSStrongauthGetU2FMetadataOutputKey `json:"content"`
}

func (s *Service) WSStrongauthGetU2FMetadata(
	input *WSStrongauthGetU2FMetadataInput,
) (
	*WSStrongauthGetU2FMetadataOutput,
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
		"https://ws1.dashlane.com/3/strongauth/getU2FMetadata",
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
	wsStrongauthGetU2FMetadataOutput := &WSStrongauthGetU2FMetadataOutput{}
	if err := json.NewDecoder(resp.Body).Decode(
		wsStrongauthGetU2FMetadataOutput,
	); err != nil {
		return nil, err
	}

	return wsStrongauthGetU2FMetadataOutput, nil
}
