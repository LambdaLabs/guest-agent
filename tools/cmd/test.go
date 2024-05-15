package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func maybeExit(err error) {
	printStack(err)
	if err != nil {
		os.Exit(1)
	}
}

func NewTestCmd(v *viper.Viper) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			maybeExit(cmd.PersistentFlags().Parse(args))

			if v.GetBool("show-config") {
				for key, value := range v.AllSettings() {
					fmt.Printf("%v: %v\n", key, value)
				}
				return
			}

			tester, err := NewTester(v)
			maybeExit(err)
			maybeExit(tester.Run())
		},
	}
	flags := cmd.PersistentFlags()

	flags.StringP("remote-host", "H", "", "Select the remote host to use for testing.")
	v.BindPFlag("guest_agent_remote_test_host", flags.Lookup("remote-host"))

	flags.StringP("remote-path", "p", "", "")
	v.BindPFlag("guest_agent_remote_test_path", flags.Lookup("remote-path"))

	flags.StringP("remote-user", "u", "", "")
	v.BindPFlag("guest_agent_remote_test_user", flags.Lookup("remote-user"))

	flags.String("debian-path", "", "Path of the guest-agent debian file to test.")
	v.BindPFlag("guest_agent_debian_path", flags.Lookup("debian-path"))

	flags.Bool("show-config", false, "show configuration")
	v.BindPFlag("show-config", flags.Lookup("show-config"))

	flags.String("expected-version", "", "set the expected version of the installed package")
	v.BindPFlag("expected-version", flags.Lookup("expected-version"))

	return cmd, nil

}

func NewTester(v *viper.Viper) (*Tester, error) {
	t := &Tester{}
	if err := v.Unmarshal(t); err != nil {
		return nil, errors.New(err)
	}
	fmt.Println(t.DebPath)
	fmt.Println(v.GetString("debian-path"))
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(t); err != nil {
		logger.Err(err).Msg("config validation failed")
		return nil, errors.New(err)
	}
	return t, nil
}

type Tester struct {
	DebPath         string `mapstructure:"guest_agent_debian_path" validate:"required"`
	ExpectedVersion string `mapstructure:"expected-version" validate:"required"`
	RemoteHost      string `mapstructure:"guest_agent_remote_test_host" validate:"required"`
	RemotePath      string `mapstructure:"guest_agent_remote_test_path" validate:"required"`
	RemoteUser      string `mapstructure:"guest_agent_remote_test_user" validate:"required"`
}

func run(command string, args ...string) error {
	logger.Info().Str("command", fmt.Sprintf("%s %s", command, strings.Join(args, " "))).Msg("running command")
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		logger.Err(err).Msg("command failed to run")
	}
	return err
}

func (t *Tester) Run() error {
	remote := fmt.Sprintf("%s@%s", t.RemoteUser, t.RemoteHost)
	if err := run("ssh", remote, "mkdir", "-p", t.RemotePath); err != nil {
		return err
	}
	debPathDest := fmt.Sprintf("%s/%s", t.RemotePath, filepath.Base(t.DebPath))
	defer func() {
		run("ssh", remote, "rm", "-rf", t.RemotePath)
	}()

	if err := run("scp", t.DebPath, fmt.Sprintf("%s:%s", remote, debPathDest)); err != nil {
		return err
	}

	testScriptSrc := "./scripts/test_script.sh"
	testScriptDest := fmt.Sprintf("%s/%s", t.RemotePath, filepath.Base(testScriptSrc))

	if err := run("scp", testScriptSrc, fmt.Sprintf("%s:%s", remote, testScriptDest)); err != nil {
		return err
	}
	if err := run("ssh", remote, "sudo", "/bin/bash", testScriptDest, debPathDest, t.ExpectedVersion); err != nil {
		return err
	}

	return nil
}
