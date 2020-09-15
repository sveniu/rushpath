package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

type APIBackupLockInput struct {
	Login *string `url:"login"`
	UKI   *string `url:"uki"`
	OTP   *string `url:"otp"`
	Lock  *string `url:"lock"`
}

// FIXME handle dynomic content
type APIBackupLockOutput struct {
	ObjectType *string `json:"objectType"`
	Content    *string `json:"content"`
}

func (s *Service) APIBackupLock(
	input *APIBackupLockInput,
) (
	*APIBackupLockOutput,
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
		"https://ws1.dashlane.com/12/backup/lock",
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
	apiBackupLockOutput := &APIBackupLockOutput{}
	if err := json.NewDecoder(resp.Body).Decode(
		apiBackupLockOutput,
	); err != nil {
		return nil, err
	}

	return apiBackupLockOutput, nil
}
