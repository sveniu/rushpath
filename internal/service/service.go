package service

const (
	// Extracted from the Firefox extension.
	defaultAppAccessKey = "C4F8H4SEAMXNBQVSASVBWDDZNCVTESMY"
	defaultAppSecretKey = "Na9Dz3WcmjMZ5pdYU1AmC5TdYkeWAOzvOK6PkbU4QjfjPQTSaXY8pjPwrvHfVH14"
)

func New() *Service {
	svc := &Service{}

	// Configure default app keys.
	appAccessKey := defaultAppAccessKey
	appSecretKey := defaultAppSecretKey
	svc.Credentials = &Credentials{}
	svc.Credentials.AppAccessKey = &appAccessKey
	svc.Credentials.AppSecretKey = &appSecretKey

	return svc
}
