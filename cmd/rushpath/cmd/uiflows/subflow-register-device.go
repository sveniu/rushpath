package uiflows

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/rs/zerolog/log"
	"github.com/sveniu/rushpath/internal/service"
)

func (ui *UI) subflowRegisterDevice() error {
	// Request registration of new device.
	apiRequestDeviceRegistrationOutput, err := ui.Service.APIRequestDeviceRegistration(
		&service.APIRequestDeviceRegistrationInput{
			Login: ui.Service.Credentials.Login,
		},
	)
	if err != nil {
		return err
	}

	verifications := apiRequestDeviceRegistrationOutput.Data.Verification
	var verificationUsable *service.APIRequestDeviceRegistrationOutputDataVerification

	for _, verification := range verifications {
		switch *verification.Type {
		case "email_token":
			if verificationUsable == nil {
				verificationUsable = verification
			}
		case "totp":
			// Skip the nil check, effectively preferring TOTP.
			verificationUsable = verification
		case "u2f":
			fmt.Printf(`
NOTE: While you have enabled U2F authentication, this program doesn't yet
support using it for logging in. In other words, the only currently supported
authentication factors are 6-digit tokens via email or TOTP.
`)
		default:
			log.Debug().
				Str("verification_type", *verification.Type).
				Msg("unknown verification type")
		}
	}

	verificationType := verificationUsable.Type
	switch *verificationType {
	case "email_token":
		// Read email token from terminal.
		emailToken, err := promptEmailToken()
		if err != nil {
			log.Error().
				Err(err).
				Msg("prompt failed")
			return err
		}

		// Request registration of new device.
		apiAuthenticationCompleteDeviceRegistrationWithTokenOutput, err := ui.Service.APIAuthenticationCompleteDeviceRegistrationWithToken(
			&service.APIAuthenticationCompleteDeviceRegistrationWithTokenInput{
				Login: ui.Service.Credentials.Login,
				Device: &service.APIAuthenticationCompleteDeviceRegistrationWithTokenInputDevice{
					DeviceName: aws.String("Firefox - Linux"),
					AppVersion: aws.String("6.2030.4"),
					Platform:   aws.String("server_standalone"),
					OSCountry:  aws.String("NO"),
					OSLanguage: aws.String("en-US"),
					Temporary:  aws.Bool(true),
				},
				Verification: &service.APIAuthenticationCompleteDeviceRegistrationWithTokenInputVerification{
					Token: &emailToken,
				},
			},
		)
		if err != nil {
			return err
		}

		// Update credentials with device keys.
		ui.Service.Credentials.DeviceAccessKey = apiAuthenticationCompleteDeviceRegistrationWithTokenOutput.Data.DeviceAccessKey
		ui.Service.Credentials.DeviceSecretKey = apiAuthenticationCompleteDeviceRegistrationWithTokenOutput.Data.DeviceSecretKey

	case "totp":
		// Read TOTP token from terminal.
		totpToken, err := promptTOTPToken(ui.Service.Credentials.OTP)
		if err != nil {
			log.Error().
				Err(err).
				Msg("prompt failed")
			return err
		}

		// Save the token to be able to check for re-use.
		ui.Service.Credentials.OTP = &totpToken

		// Request registration of new device.
		apiAuthenticationCompleteDeviceRegistrationWithTOTPOutput, err := ui.Service.APIAuthenticationCompleteDeviceRegistrationWithTOTP(
			&service.APIAuthenticationCompleteDeviceRegistrationWithTOTPInput{
				Login: ui.Service.Credentials.Login,
				Device: &service.APIAuthenticationCompleteDeviceRegistrationWithTOTPInputDevice{
					DeviceName: aws.String("Firefox - Linux"),
					AppVersion: aws.String("6.2030.4"),
					Platform:   aws.String("server_standalone"),
					OSCountry:  aws.String("NO"),
					OSLanguage: aws.String("en-US"),
					Temporary:  aws.Bool(true),
				},
				Verification: &service.APIAuthenticationCompleteDeviceRegistrationWithTOTPInputVerification{
					OTP: &totpToken,
				},
			},
		)
		if err != nil {
			return err
		}

		// Update credentials with device keys.
		ui.Service.Credentials.DeviceAccessKey = apiAuthenticationCompleteDeviceRegistrationWithTOTPOutput.Data.DeviceAccessKey
		ui.Service.Credentials.DeviceSecretKey = apiAuthenticationCompleteDeviceRegistrationWithTOTPOutput.Data.DeviceSecretKey

	case "u2f":
		err = fmt.Errorf("u2f auth not yet implemented")
		return err

	default:
		return fmt.Errorf(
			"unknown auth type `%s`",
			*verificationType,
		)
	}

	// The UKI is a simple concatenation of device keys.
	uki := fmt.Sprintf(
		"%s-%s",
		*ui.Service.Credentials.DeviceAccessKey,
		*ui.Service.Credentials.DeviceSecretKey,
	)
	ui.Service.Credentials.UKI = &uki

	return nil
}
