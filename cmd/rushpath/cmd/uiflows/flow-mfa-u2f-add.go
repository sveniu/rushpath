package uiflows

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func (ui *UI) MFAU2FAdd() error {
	fmt.Print(`
Adding authentication factor: Universal 2nd Factor (U2F)

	`)

	if err := ui.loginAndRegisterDevice(); err != nil {
		return err
	}

	if err := ui.subflowRegisterU2F(); err != nil {
		log.Error().
			Err(err).
			Msg("u2f registration failed")
		return err
	}

	if err := ui.subflowU2FListKeys(); err != nil {
		log.Error().
			Err(err).
			Msg("u2f key listing failed")
		return err
	}

	fmt.Printf("\nU2F device successfully added!\n\n")

	return nil
}
