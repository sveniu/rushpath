package uiflows

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sveniu/rushpath/internal/service"
)

func (ui *UI) subflowSetRecoveryPhone() error {
	recoveryPhone, err := promptPhone()
	if err != nil {
		log.Error().
			Err(err).
			Msg("phone prompt failed")
		return err
	}

	if recoveryPhone == "" {
		fmt.Printf("\nSkipping recovery phone setup.\n")
		return nil
	}

	_, err = ui.Service.WSStrongauthSetRecoveryPhone(
		&service.WSStrongauthSetRecoveryPhoneInput{
			Login: ui.Service.Credentials.Login,
			UKI:   ui.Service.Credentials.UKI,
			Phone: &recoveryPhone,
		},
	)
	if err != nil {
		return err
	}

	log.Info().
		Str("phone", recoveryPhone).
		Msg("recovery phone set")

	return nil
}
