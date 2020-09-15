package uiflows

import (
	"fmt"
	"time"

	"github.com/hako/durafmt"
	"github.com/rs/zerolog/log"
	"github.com/sveniu/rushpath/internal/service"
)

func (ui *UI) subflowU2FListKeys() error {
	wsStrongauthGetU2FMetadataOutput, err := ui.Service.WSStrongauthGetU2FMetadata(
		&service.WSStrongauthGetU2FMetadataInput{
			Login: ui.Service.Credentials.Login,
			UKI:   ui.Service.Credentials.UKI,
		},
	)
	if err != nil {
		return err
	}

	if len(wsStrongauthGetU2FMetadataOutput.Content) < 1 {
		log.Info().
			Msg("no u2f keys registered")
		return nil
	}

	fmt.Printf(`
List of registered U2F keys:

`)

	for _, key := range wsStrongauthGetU2FMetadataOutput.Content {
		creationDate := time.Unix(*key.CreationDateUnix, 0)
		creationDateStr := creationDate.UTC().Format("2006-01-02 15:04:05Z")
		creationDateAge := durafmt.Parse(time.Now().Sub(creationDate)).LimitFirstN(2)

		lastUsedDate := time.Unix(*key.LastUsedDateUnix, 0)
		lastUsedDateStr := lastUsedDate.UTC().Format("2006-01-02 15:04:05Z")
		lastUsedDateAge := durafmt.Parse(time.Now().Sub(lastUsedDate)).LimitFirstN(2)

		currentIPStr := "note: this matches your current IP"
		if !*key.LastUsedFromThisIP {
			currentIPStr = "note: this does not match your current IP"
		}

		fmt.Printf("  Key name:       %s\n", *key.Name)
		fmt.Printf("  Key handle:     %s\n", *key.KeyHandle)
		fmt.Printf("  Creation date:  %s (age %s)\n", creationDateStr, creationDateAge)
		fmt.Printf("  Last used:\n")
		fmt.Printf("    On date:      %s (age %s)\n", lastUsedDateStr, lastUsedDateAge)
		fmt.Printf("    From IP:      %s (%s)\n", *key.LastUsedFromIP, currentIPStr)
		fmt.Printf("    From country: %s\n", *key.LastUsedFromCountry)
		fmt.Printf("\n")
	}

	return nil
}
