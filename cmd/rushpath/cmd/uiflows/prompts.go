package uiflows

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/rs/zerolog/log"
)

func promptLogin() (login string, err error) {
	fmt.Print(`
Please enter your Dashlane login.

	`)

	loginPrompt := &survey.Input{
		Message: "Login:",
	}
	if err = survey.AskOne(loginPrompt, &login); err != nil {
		log.Error().
			Err(err).
			Msg("login prompt failed")
		return
	}

	return
}

func promptPhone() (phone string, err error) {
	fmt.Print(`
Please enter your recovery phone number.

	`)

	phonePrompt := &survey.Input{
		Message: "Recovery phone:",
	}
	if err = survey.AskOne(phonePrompt, &phone); err != nil {
		log.Error().
			Err(err).
			Msg("login prompt failed")
		return
	}

	return
}

func promptEmailToken() (emailToken string, err error) {
	fmt.Print(`
Please enter the 6-digit code received by email.

	`)

	emailTokenPrompt := &survey.Input{
		Message: "Email token:",
	}
	if err = survey.AskOne(
		emailTokenPrompt,
		&emailToken,
		survey.WithValidator(func(v interface{}) error {
			if len(v.(string)) != 6 {
				return fmt.Errorf("Token must be 6 digits")
			}
			return nil
		}),
	); err != nil {
		log.Error().
			Err(err).
			Msg("email token prompt failed")
		return
	}

	return
}

func promptTOTPToken(oldTOTPToken *string) (totpToken string, err error) {
	fmt.Print(`
Please enter the 6-digit TOTP code from your auth app.

	`)

	totpTokenPrompt := &survey.Input{
		Message: "TOTP token:",
	}
	if err = survey.AskOne(
		totpTokenPrompt,
		&totpToken,
		survey.WithValidator(func(v interface{}) error {
			if len(v.(string)) != 6 {
				return fmt.Errorf("Token must be 6 digits")
			}
			return nil
		}),
		survey.WithValidator(func(ans interface{}) error {
			if oldTOTPToken != nil && ans.(string) == *oldTOTPToken {
				return fmt.Errorf("Token already used; please try a new one")
			}
			return nil
		}),
	); err != nil {
		log.Error().
			Err(err).
			Msg("totp token prompt failed")
		return
	}

	return
}

func promptU2FDeviceName() (login string, err error) {
	fmt.Print(`
Please enter the name of your U2F device.

	`)

	loginPrompt := &survey.Input{
		Message: "U2F device name:",
	}
	if err = survey.AskOne(loginPrompt, &login); err != nil {
		log.Error().
			Err(err).
			Msg("u2f device name prompt failed")
		return
	}

	return
}
