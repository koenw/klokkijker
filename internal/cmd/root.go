package cmd

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	rootCmd = &cobra.Command{
		Use:  "klokkijker <NTP servers>",
		Long: `Diagnostic NTP command-line client & prometheus metrics exporter.`,
		Example: `  klokkijker ntp.example.com
  klokkijker ntp1.example.com ntp2.example.com

  klokkijker ping --format=json ntp.example.com |jq .fields.offset
  klokkijker ping --interval 1 --count 3 ntp.example.com

  klokkijker monitor ntp.example.com
  klokkijker monitor --interval 1 --prometheus-address 0.0.0.0 ntp1.example.com ntp2.example.com

  klokkijker loadgen --rpm 10000 --ramp-period 5m ntp.example.com
  klokkijker loadgen --rpm 1000 --ramp-period 120s ntp1.example.com ntp2.example.com`,
	}
	defaultCmd   = pingCmd
	outputFormat string
)

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.PersistentFlags().StringVarP(
		&outputFormat,
		"format", "f",
		"cli",
		"Output format, one of 'json', 'cli' or 'text'")
}

func subCommands() (commandNames []string) {
	for _, command := range rootCmd.Commands() {
		commandNames = append(commandNames, append(command.Aliases, command.Name())...)
	}
	commandNames = append(commandNames, "help")
	commandNames = append(commandNames, "completion")
	return
}

func setDefaultCommandIfNonePresent(defaultCommand string) {
	if len(os.Args) > 1 {
		potentialCommand := os.Args[1]
		for _, command := range subCommands() {
			if command == potentialCommand {
				return
			}
		}
		if rootCmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
			os.Args = append([]string{os.Args[0], defaultCommand}, os.Args[1:]...)
		}
	}
}

func setupLogging(format string) {
	switch format {
	case "json":
		log.SetHandler(json.New(os.Stdout))
	case "text":
		log.SetHandler(text.New(os.Stdout))
	default:
		log.SetHandler(cli.New(os.Stderr))
	}
}

func Execute() {
	setDefaultCommandIfNonePresent("ping")

	switch outputFormat {
	case "json":
		log.SetHandler(json.New(os.Stdout))
	case "text":
		log.SetHandler(text.New(os.Stdout))
	default:
		log.SetHandler(cli.New(os.Stderr))
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
