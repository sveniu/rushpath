package uiflows

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func (ui *UI) loginAndRegisterDevice() error {
	login, err := promptLogin()
	if err != nil {
		return err
	}

	ui.Service.Credentials.Login = &login

	if err := ui.subflowLoginExists(); err != nil {
		log.Error().
			Err(err).
			Msg("login check failed")
		return err
	}

	log.Info().
		Str("login", login).
		Msg("login exists")

	if err := ui.subflowRegisterDevice(); err != nil {
		log.Error().
			Err(err).
			Msg("device registration failed")
		return err
	}

	log.Info().
		Msg("device registered")

	fmt.Printf(`
Successfully registered temporary device, which will be used to perform the
desired action.
`)

	return nil
}
