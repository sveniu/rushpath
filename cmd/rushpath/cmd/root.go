package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var cmdRoot = &cobra.Command{
	SilenceUsage: true,
	Use:          "rushpath",
	Short:        "Rushpath – a Dashlane CLI",
	Long: `Rushpath is a simple CLI for the Dashlane password manager.

The basic idea is to allow for TOTP and U2F management to platforms other
than the Windows and OSX desktop applications.

THERE ARE KNOWN BUGS, SO USE AT YOUR OWN RISK!

https://github.com/sveniu/rushpath
`,
}

func init() {
	cobra.OnInitialize(initLogging)

	cmdRoot.PersistentFlags().StringP(
		"loglevel",
		"l",
		"warn",
		`log level: error, warn, info, debug or trace
`,
	)
}

func initLogging() {
	// Log to the terminal.
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		},
	)

	// Get log level from flag.
	logLevelString, err := cmdRoot.PersistentFlags().GetString("loglevel")
	if err != nil {
		log.Warn().
			Err(err).
			Msg("error reading log level")
		return
	}

	logLevel, err := zerolog.ParseLevel(logLevelString)
	if err != nil {
		logLevel = zerolog.WarnLevel
		log.Warn().
			Err(err).
			Str("log_level", logLevel.String()).
			Msg("revert to default log level")
	}

	zerolog.SetGlobalLevel(logLevel)

	log.Info().
		Str("log_level", log.Logger.GetLevel().String()).
		Msg("log level configured")
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once to
// the cmdRoot.
func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
