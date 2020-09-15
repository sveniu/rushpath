package uiflows

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func (ui *UI) MFAU2FRemove() error {
	fmt.Print(`
Removing authentication factor: Universal 2nd Factor (U2F)

	`)

	if err := ui.loginAndRegisterDevice(); err != nil {
		return err
	}

	if err := ui.subflowDeregisterU2F(); err != nil {
		log.Error().
			Err(err).
			Msg("u2f deregistration failed")
		return err
	}

	fmt.Printf("\nU2F device(s) successfully removed!\n\n")

	return nil
}
