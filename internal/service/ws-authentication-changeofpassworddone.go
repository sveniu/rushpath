package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

type WSAuthenticationChangeOfPasswordDoneInput struct {
	Login *string `url:"login"`
	UKI   *string `url:"uki"`
}

type WSAuthenticationChangeOfPasswordDoneOutput struct {
	Result *string `json:"result"`
}

func (s *Service) WSAuthenticationChangeOfPasswordDone(
	input *WSAuthenticationChangeOfPasswordDoneInput,
) (
	*WSAuthenticationChangeOfPasswordDoneOutput,
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
		"https://ws1.dashlane.com/7/authentication/changeOfPasswordDone",
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
	wsAuthenticationChangeOfPasswordDoneOutput := &WSAuthenticationChangeOfPasswordDoneOutput{}
	if err := json.NewDecoder(resp.Body).Decode(
		wsAuthenticationChangeOfPasswordDoneOutput,
	); err != nil {
		return nil, err
	}

	return wsAuthenticationChangeOfPasswordDoneOutput, nil
}
