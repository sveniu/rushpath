package uiflows

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sveniu/rushpath/internal/service"
)

func (ui *UI) MFATOTPDisable() error {
	fmt.Print(`
Disabling authentication factor: Time-based One-Time Password (TOTP)

	`)

	if err := ui.loginAndRegisterDevice(); err != nil {
		return err
	}

	fmt.Printf(`
To disable TOTP, you have to supply a final TOTP code. Make sure to not
re-use the previous one. If your app still shows the previous one, just wait
until a new code appears.

`)

	totpToken, err := promptTOTPToken(ui.Service.Credentials.OTP)
	if err != nil {
		return err
	}

	// Save the token to be able to check for re-use.
	ui.Service.Credentials.OTP = &totpToken

	wsStrongauthDeactivateForNewDeviceOutput, err := ui.Service.WSStrongauthDeactivateForNewDevice(
		&service.WSStrongauthDeactivateForNewDeviceInput{
			Login: ui.Service.Credentials.Login,
			UKI:   ui.Service.Credentials.UKI,
			OTP:   ui.Service.Credentials.OTP,
		},
	)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to deactivate for new device")
	}

	if *wsStrongauthDeactivateForNewDeviceOutput.Code != 200 ||
		*wsStrongauthDeactivateForNewDeviceOutput.Message != "OK" {
		err := fmt.Errorf("non-200 or non-OK return")
		log.Error().
			Err(err).
			Msg("failed to deactivate TOTP for new devices")
		return err
	}

	log.Info().
		Msg("disabled TOTP for new devices")

	fmt.Printf("\nTOTP disabled successfully!\n\n")

	return nil
}
