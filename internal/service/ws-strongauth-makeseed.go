package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

type WSStrongauthMakeSeedInput struct {
	Login *string `url:"login"`
	UKI   *string `url:"uki"`
}

type WSStrongauthMakeSeedOutputContentQRCode struct {
	Data  *string `json:"data"`
	Width *int    `json:"width"`
}

type WSStrongauthMakeSeedOutputContent struct {
	Seed      *string                                  `json:"seed"`
	ServerKey *string                                  `json:"serverKey"`
	QRCode    *WSStrongauthMakeSeedOutputContentQRCode `json:"qrcode"`
}

type WSStrongauthMakeSeedOutput struct {
	Code    *int                               `json:"code"`
	Message *string                            `json:"message"`
	Content *WSStrongauthMakeSeedOutputContent `json:"content"`
}

type WSStrongauthMakeSeedError struct {
	Code    *int    `json:"code"`
	Message *string `json:"message"`
	Content struct {
		Status *string `json:"status"`
	} `json:"content"`
}

type ErrTOTPAlreadyEnabled struct{}

func (e *ErrTOTPAlreadyEnabled) Error() string {
	return "totp already enabled"
}

func (s *Service) WSStrongauthMakeSeed(
	input *WSStrongauthMakeSeedInput,
) (
	*WSStrongauthMakeSeedOutput,
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
		"https://ws1.dashlane.com/3/strongauth/makeSeed",
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

	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Decode the response. Try errors first.
	wsStrongauthMakeSeedError := &WSStrongauthMakeSeedError{}
	if err := json.NewDecoder(bytes.NewReader(responseBodyBytes)).Decode(
		wsStrongauthMakeSeedError,
	); err != nil {
		return nil, err
	}

	/*
		{
			"code": 403,
			"message": "Forbidden",
			"content": {
				"status": "newDevice"
			}
		}
	*/
	if *wsStrongauthMakeSeedError.Code == 403 &&
		wsStrongauthMakeSeedError.Content.Status != nil &&
		*wsStrongauthMakeSeedError.Content.Status == "newDevice" {
		// FIXME what about non-"newDevice" - how does it look?
		return nil, &ErrTOTPAlreadyEnabled{}
	}

	// No errors found; decode as successful response.
	wsStrongauthMakeSeedOutput := &WSStrongauthMakeSeedOutput{}
	if err := json.NewDecoder(bytes.NewReader(responseBodyBytes)).Decode(
		wsStrongauthMakeSeedOutput,
	); err != nil {
		return nil, err
	}

	return wsStrongauthMakeSeedOutput, nil
}
