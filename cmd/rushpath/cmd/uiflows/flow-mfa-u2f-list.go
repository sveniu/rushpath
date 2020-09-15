package uiflows

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func (ui *UI) MFAU2FList() error {
	fmt.Print(`
Listing authentication factors: Universal 2nd Factor (U2F)

	`)

	if err := ui.loginAndRegisterDevice(); err != nil {
		return err
	}

	if err := ui.subflowU2FListKeys(); err != nil {
		log.Error().
			Err(err).
			Msg("u2f key listing failed")
		return err
	}

	fmt.Printf("\nU2F devices listed successfully!\n\n")

	return nil
}
