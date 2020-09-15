package uiflows

import (
	"encoding/base64"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/sveniu/rushpath/internal/service"
)

func (ui *UI) subflowRegisterU2F() error {
	wsStrongauthStartU2FRegistrationOutput, err := ui.Service.WSStrongauthStartU2FRegistration(
		&service.WSStrongauthStartU2FRegistrationInput{
			Login: ui.Service.Credentials.Login,
			UKI:   ui.Service.Credentials.UKI,
		},
	)
	if err != nil {
		return err
	}

	// Construct the client data.
	// https://fidoalliance.org/specs/fido-u2f-v1.2-ps-20170411/fido-u2f-raw-message-formats-v1.2-ps-20170411.html#idl-def-ClientData
	clientDataTyp := "navigator.id.finishEnrollment"
	clientDataJSON, err := json.Marshal(struct {
		Challenge *string `json:"challenge"`
		Origin    *string `json:"origin"`
		Typ       *string `json:"typ"`
	}{
		Challenge: wsStrongauthStartU2FRegistrationOutput.Content.Challenge,
		Origin:    wsStrongauthStartU2FRegistrationOutput.Content.AppID,
		Typ:       &clientDataTyp,
	})
	if err != nil {
		return err
	}

	regData, err := promptU2FRegister(
		clientDataJSON,
		[]byte(*wsStrongauthStartU2FRegistrationOutput.Content.AppID),
	)
	if err != nil {
		return err
	}

	registrationDataB64 := base64.RawURLEncoding.EncodeToString(regData)
	clientDataB64 := base64.RawURLEncoding.EncodeToString(clientDataJSON)

	challengeAnswerJSON, err := json.Marshal(struct {
		RegistrationData *string `json:"registrationData"`
		ClientData       *string `json:"clientData"`
	}{
		RegistrationData: &registrationDataB64,
		ClientData:       &clientDataB64,
	})
	if err != nil {
		return err
	}

	challengeAnswerJSONString := string(challengeAnswerJSON)

	log.Trace().
		Bytes("client_data_json", clientDataJSON).
		Str("challenge_answer_json", challengeAnswerJSONString).
		Msg("generated challenge answer")

	deviceName, err := promptU2FDeviceName()
	if err != nil {
		return err
	}

	wsStrongauthRegisterU2FDeviceOutput, err := ui.Service.WSStrongauthRegisterU2FDevice(
		&service.WSStrongauthRegisterU2FDeviceInput{
			ChallengeAnswer: &challengeAnswerJSONString,
			Name:            &deviceName,
			Login:           ui.Service.Credentials.Login,
			UKI:             ui.Service.Credentials.UKI,
		},
	)
	if err != nil {
		return err
	}

	log.Debug().
		Int("code", *wsStrongauthRegisterU2FDeviceOutput.Code).
		Str("message", *wsStrongauthRegisterU2FDeviceOutput.Message).
		Msg("u2f dev registration done?")

	return nil
}
