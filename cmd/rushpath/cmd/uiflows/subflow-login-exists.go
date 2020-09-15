package uiflows

import (
	"fmt"

	"github.com/sveniu/rushpath/internal/service"
)

func (ui *UI) subflowLoginExists() error {
	wsAuthenticationExistsOutput, err := ui.Service.WSAuthenticationExists(
		&service.WSAuthenticationExistsInput{
			Login: ui.Service.Credentials.Login,
		},
	)
	if err != nil {
		return err
	}

	if wsAuthenticationExistsOutput.Exists == nil ||
		(*wsAuthenticationExistsOutput.Exists != "YES" &&
			*wsAuthenticationExistsOutput.Exists != "YES_OTP_NEWDEVICE") {
		return fmt.Errorf(
			"Login does not exist: %s",
			*ui.Service.Credentials.Login,
		)
	}

	return nil
}
