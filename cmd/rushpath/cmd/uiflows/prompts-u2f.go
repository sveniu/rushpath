package uiflows

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/flynn/u2f/u2fhid"
	"github.com/flynn/u2f/u2ftoken"
	"github.com/rs/zerolog/log"
)

func promptU2FRegister(clientData, appID []byte) (regData []byte, err error) {
	deviceInfos, err := u2fhid.Devices()
	if err != nil {
		return
	}
	if len(deviceInfos) == 0 {
		err = fmt.Errorf("no U2F devices found")
		return
	}

	fmt.Printf("\nFound U2F device")
	if len(deviceInfos) > 1 {
		fmt.Printf("s")
	}
	fmt.Printf(":\n\n")

	// List of U2F devices for later selection display and indexing.
	devicesDisplay := []string{}

	for _, device := range deviceInfos {
		log.Debug().
			Str("vendor_id", fmt.Sprintf("0x%04x", device.VendorID)).
			Str("product_id", fmt.Sprintf("0x%04x", device.ProductID)).
			Str("version_number", fmt.Sprintf("0x%04x", device.VersionNumber)).
			Str("manufacturer", device.Manufacturer).
			Str("product", device.Product).
			Msg("found u2f device")

		deviceDisplay := fmt.Sprintf(
			"%s %s — vendor_id=%s product_id=%s version=%s",
			device.Manufacturer,
			device.Product,
			fmt.Sprintf("0x%04x", device.VendorID),
			fmt.Sprintf("0x%04x", device.ProductID),
			fmt.Sprintf("0x%04x", device.VersionNumber),
		)

		fmt.Printf("  %s\n", deviceDisplay)

		// Append to the display list.
		devicesDisplay = append(devicesDisplay, deviceDisplay)
	}
	fmt.Println()

	// Default to the first device.
	deviceIndex := 0

	// Prompt if there are more devices.
	if len(deviceInfos) > 1 {
		prompt := &survey.Select{
			Message: "Select which U2F device to use:",
			Options: devicesDisplay,
		}
		survey.AskOne(prompt, &deviceIndex)
	}

	device, err := u2fhid.Open(deviceInfos[deviceIndex])
	if err != nil {
		return
	}
	token := u2ftoken.NewToken(device)

	version, err := token.Version()
	if err != nil {
		return
	}
	log.Debug().
		Str("u2f_version", version).
		Msg("got token u2f version")

	fmt.Printf("Touch the gold disc of your U2F device to continue.\n")

	for {
		regData, err = token.Register(
			u2ftoken.RegisterRequest{
				Challenge:   hashSHA256(clientData),
				Application: hashSHA256(appID),
			},
		)
		if err == u2ftoken.ErrPresenceRequired {
			time.Sleep(200 * time.Millisecond)
			continue
		} else if err != nil {
			return
		}
		break
	}

	log.Debug().
		Hex("registration_response_hex", regData).
		Msg("registered device")

	device.Close()

	return
}

func hashSHA256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}
