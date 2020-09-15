package service

import "net/http"

type Credentials struct {
	Login           *string `json:"login,omitempty"`
	AppAccessKey    *string `json:"app_access_key,omitempty"`
	AppSecretKey    *string `json:"app_secret_key,omitempty"`
	DeviceAccessKey *string `json:"device_access_key,omitempty"`
	DeviceSecretKey *string `json:"device_secret_key,omitempty"`
	UKI             *string `json:"uki,omitempty"`
	OTP             *string `json:"otp,omitempty"`
	LockID          *string `json:"lock_id,omitempty"`
}

type Service struct {
	Client      *http.Client `json:"-"`
	Credentials *Credentials `json:"credentials"`
}

type authType int

const (
	authTypeNone = iota
	authTypeApp
	authTypeUserDevice
	authTypeSession
	authTypeTeamDevice
)
