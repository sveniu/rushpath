package uiflows

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func (ui *UI) MFATOTPEnable() error {
	fmt.Print(`
Enabling authentication factor: Time-based One-Time Password (TOTP)

	`)

	if err := ui.loginAndRegisterDevice(); err != nil {
		return err
	}

	if err := ui.subflowRegisterTOTP(); err != nil {
		log.Error().
			Err(err).
			Msg("totp registration failed")
		return err
	}

	log.Info().
		Msg("totp registered")

	fmt.Print(`
Provide a backup phone number in case you lose your primary phone. Dashlane
will only use this number to text (SMS) you a security code.

Note: An empty value will skip this step.

	`)

	if err := ui.subflowSetRecoveryPhone(); err != nil {
		log.Error().
			Err(err).
			Msg("set recovery phone failed")
		return err
	}

	fmt.Printf("\nTOTP enabled successfully!\n\n")

	return nil
}
