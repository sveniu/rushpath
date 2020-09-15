package uiflows

import (
	"fmt"
	"os"
	"strings"

	"github.com/sveniu/rushpath/internal/service"

	"github.com/AlecAivazis/survey/v2"
	"github.com/google/uuid"
	"github.com/mdp/qrterminal"
	"github.com/rs/zerolog/log"
)

func (ui *UI) subflowRegisterTOTP() error {
	// Start 2FA registration.
	wsStrongauthMakeSeedOutput, err := ui.Service.WSStrongauthMakeSeed(
		&service.WSStrongauthMakeSeedInput{
			Login: ui.Service.Credentials.Login,
			UKI:   ui.Service.Credentials.UKI,
		},
	)
	if err != nil {
		if _, ok := err.(*service.ErrTOTPAlreadyEnabled); ok {
			// FIXME show if it's for 'new devices' only, or for every login.
			fmt.Printf(`
TOTP is already enabled. You can switch the TOTP prompt type, i.e. whether to
require TOTP for new devices only or for every login, by disabling TOTP and
re-enabling it again.
`)
		}
		return err
	}

	// Prepare and print QR code.
	qrData := fmt.Sprintf(
		"otpauth://totp/Dashlane:%s?secret=%s&issuer=Dashlane",
		*ui.Service.Credentials.Login,
		*wsStrongauthMakeSeedOutput.Content.Seed,
	)

	fmt.Printf(`
A large ~90x50 character QR code is about to be printed on your terminal. 

Note: This prompt is merely for you to choose whether to hide or obscure the
screen from prying eyes before printing the code. If you answer no on the
prompt, the entire TOTP registration is with an error, and you would have to
start over.

`)

	qrConfirm := false
	qrConfirmPrompt := &survey.Confirm{
		Message: "Display QR code (sensitive data)?",
		Default: true,
	}
	if err := survey.AskOne(qrConfirmPrompt, &qrConfirm); err != nil {
		return err
	}

	if !qrConfirm {
		return fmt.Errorf("cancelling on user request")
	}

	qrterminal.Generate(qrData, qrterminal.L, os.Stdout)
	fmt.Println()

	totpKeyB32 := *wsStrongauthMakeSeedOutput.Content.Seed
	totpKeyB32Formatted := fmt.Sprintf(
		"%s-%s-%s-%s %s-%s-%s-%s\n",
		totpKeyB32[0:4],
		totpKeyB32[4:8],
		totpKeyB32[8:12],
		totpKeyB32[12:16],
		totpKeyB32[16:20],
		totpKeyB32[20:24],
		totpKeyB32[24:28],
		totpKeyB32[28:32],
	)
	fmt.Printf(
		"Key for manual entry (ignore dashes and spaces): %s\n",
		totpKeyB32Formatted,
	)

	fmt.Printf(`
To continue, scan the QR code or enter the manual code (without dashes and
spaces). Then type in the 6-digit code from the authentication app in the
prompt below.
`)
	totpToken, err := promptTOTPToken(ui.Service.Credentials.OTP)
	if err != nil {
		log.Error().
			Err(err).
			Msg("prompt failed")
		return err
	}

	// Save the token to be able to check for re-use.
	ui.Service.Credentials.OTP = &totpToken

	// FIXME why is offline access impossible?
	fmt.Print(`
TOTP can be enabled in two ways: Only for new devices, or for every login.

If enabled for every login, offline access to Dashlane will not be possible.

`)

	// FIXME move to prompts.go
	var totpSchemeIndex int
	totpSchemePrompt := &survey.Select{
		Message: "Enable TOTP MFA",
		Options: []string{
			"Only for new device logins",
			"For every login",
		},
	}
	if err := survey.AskOne(totpSchemePrompt, &totpSchemeIndex); err != nil {
		log.Error().
			Err(err).
			Msg("prompt failed")
		return err
	}

	switch totpSchemeIndex {
	case 0: // new devices
		wsStrongauthActivateForNewDeviceOutput, err := ui.Service.WSStrongauthActivateForNewDevice(
			&service.WSStrongauthActivateForNewDeviceInput{
				Login: ui.Service.Credentials.Login,
				UKI:   ui.Service.Credentials.UKI,
				OTP:   ui.Service.Credentials.OTP,
			},
		)
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to activate for new device")
		}

		if *wsStrongauthActivateForNewDeviceOutput.Code != 200 ||
			*wsStrongauthActivateForNewDeviceOutput.Message != "OK" {
			err := fmt.Errorf("non-200 or non-OK return")
			log.Error().
				Err(err).
				Msg("failed to activate TOTP for new devices")
			return err
		}

		log.Info().
			Msg("enabled TOTP for new devices")

	case 1: // always
		// Generate lock UUID.
		lockUUID, err := uuid.NewRandom()
		if err != nil {
			return nil
		}

		// Lock format is {UUID}, but maybe it's not required?
		lockID := fmt.Sprintf(
			"{%s}",
			strings.ToUpper(lockUUID.String()),
		)
		ui.Service.Credentials.LockID = &lockID

		// Start 2FA registration.
		// FIXME this doesn't work; rather use backup/upload with
		// form item strongAuthSetting=login.
		apiBackupLockOutput, err := ui.Service.APIBackupLock(
			&service.APIBackupLockInput{
				Login: ui.Service.Credentials.Login,
				UKI:   ui.Service.Credentials.UKI,
				OTP:   ui.Service.Credentials.OTP,
				Lock:  ui.Service.Credentials.LockID,
			},
		)
		if err != nil {
			return err
		}

		log.Debug().
			Str("backup_lock_output", fmt.Sprintf("%v", *apiBackupLockOutput)).
			Msg("got backup lock response")

	default:
		err := fmt.Errorf("invalid choice")
		log.Error().
			Err(err).
			Msg("invalid choice")
		return err
	}

	// Get recovery keys.
	wsStrongauthGetRecoveryKeysOutput, err := ui.Service.WSStrongauthGetRecoveryKeys(
		&service.WSStrongauthGetRecoveryKeysInput{
			Login: ui.Service.Credentials.Login,
			UKI:   ui.Service.Credentials.UKI,
		},
	)
	if err != nil {
		return err
	}

	recoveryCodeConfirm := false
	recoveryCodeConfirmPrompt := &survey.Confirm{
		Message: "Display recovery codes (sensitive data)?",
		Default: true,
	}
	if err := survey.AskOne(
		recoveryCodeConfirmPrompt,
		&recoveryCodeConfirm,
	); err != nil {
		return err
	}

	if !recoveryCodeConfirm {
		return fmt.Errorf("cancelling on user request")
	}

	fmt.Printf("\nRecovery keys:\n\n")
	for _, key := range wsStrongauthGetRecoveryKeysOutput.Content {
		fmt.Printf("  %s\n", *key)
	}
	fmt.Println()

	return nil
}
