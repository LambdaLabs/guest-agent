package cmd

import (
	"fmt"
	"os"

	"github.com/go-errors/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logger zerolog.Logger

type timestampHook struct{}

func (t timestampHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	e.Timestamp()
}

func NewRootCmd() *cobra.Command {
	v := viper.New()
	v.SetConfigType("env")
	v.SetConfigName("guest-agent")
	v.AddConfigPath(".")
	v.AddConfigPath("../")
	v.SetEnvPrefix("GUEST_AGENT")
	cmd := &cobra.Command{
		Use: "guest_agent [command]",
	}

	logger = zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stdout,
	}).Hook(timestampHook{})

	subCommands := []func(v *viper.Viper) (*cobra.Command, error){
		NewRenderTemplateCmd,
		NewTagCmd,
	}
	for _, CommandFunc := range subCommands {
		subCmd, err := CommandFunc(v)
		if err != nil {
			panic(err)
		}
		cmd.AddCommand(subCmd)
	}
	return cmd
}

func printStack(err error) {
	if err == nil {
		return
	}
	newErr, ok := err.(*errors.Error)
	if ok {
		fmt.Printf("%v\n", newErr.ErrorStack())
	} else {
		fmt.Printf("%v\n", err)
	}
}
