package cmd

import (
	"fmt"
	"os"

	"text/template"

	"github.com/chigopher/pathlib"
	"github.com/go-errors/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd() *cobra.Command {
	v := viper.New()
	v.SetConfigType("env")
	v.SetConfigName("guest-agent")
	v.AddConfigPath(".")
	v.SetEnvPrefix("GUEST_AGENT")
	cmd := &cobra.Command{
		Use: "guest_agent [command]",
	}
	subCmd, err := NewRenderTemplateCmd(v)
	if err != nil {
		panic(err)
	}
	cmd.AddCommand(subCmd)
	return cmd
}

func printStack(err error) {
	if err == nil {
		return
	}
	newErr, ok := err.(*errors.Error)
	if ok {
		fmt.Println(newErr.ErrorStack())
	}
}

func NewRenderTemplateCmd(v *viper.Viper) (*cobra.Command, error) {
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	cmd := &cobra.Command{
		Use: "render_template",
		RunE: func(cmd *cobra.Command, args []string) error {
			renderer, err := GetNewTemplateRendererFromViper(v)
			if err != nil {
				printStack(err)
				return err
			}
			err = renderer.Run()
			printStack(err)
			return err
		},
	}
	return cmd, nil
}

type rendererConf struct {
	EtcBaseDir                   string `mapstructure:"guest_agent_etc_base_dir"`
	TelegrafConf                 string `mapstructure:"guest_agent_telegraf_conf"`
	GuestAgentLambdaBinDir       string `mapstructure:"guest_agent_lambda_bin_dir"`
	GuestAgentServiceName        string `mapstructure:"guest_agent_service_name"`
	GuestAgentServiceFile        string `mapstructure:"guest_agent_service_file"`
	GuestAgentUpdaterServiceName string `mapstructure:"guest_agent_updater_service_name"`
	GuestAgentUpdaterServiceFile string `mapstructure:"guest_agent_updater_service_file"`
	TemplatesDir                 string `mapstructure:"guest_agent_templates_dir"`
	TemplatesOutDir              string `mapstructure:"guest_agent_templates_outdir"`
}

type TemplateRenderer struct {
	config rendererConf
}

func GetNewTemplateRendererFromViper(v *viper.Viper) (*TemplateRenderer, error) {

	rendererConf := rendererConf{}
	if err := v.Unmarshal(&rendererConf); err != nil {
		return nil, errors.Join(err)
	}
	if rendererConf.TemplatesDir == "" {
		return nil, errors.New("must specify templates dir")
	}
	return &TemplateRenderer{
		config: rendererConf,
	}, nil
}

func (r *TemplateRenderer) Run() error {
	templatesDir := pathlib.NewPath(r.config.TemplatesDir)
	outputDir := pathlib.NewPath(r.config.TemplatesOutDir)

	walker, err := pathlib.NewWalk(
		templatesDir,
		pathlib.WalkVisitDirs(false),
		pathlib.WalkVisitFiles(true),
	)
	if err != nil {
		return errors.Join(err)
	}

	err = walker.Walk(func(path *pathlib.Path, info os.FileInfo, err error) error {
		template := template.New(path.String())

		fileBytes, err := path.ReadFile()
		if err != nil {
			return errors.Join(err)
		}
		parsed, err := template.Parse(string(fileBytes))
		if err != nil {
			return errors.Join(err)
		}
		pathRelative, err := path.RelativeTo(templatesDir)
		if err != nil {
			return errors.Join(err)
		}
		fileOutpath := outputDir.JoinPath(pathRelative)
		if err := fileOutpath.Parent().MkdirAll(); err != nil {
			return errors.Join(err)
		}
		outFile, err := fileOutpath.OpenFile(os.O_WRONLY | os.O_TRUNC | os.O_CREATE)
		if err != nil {
			return errors.Join(err)
		}
		defer outFile.Close()

		if err := parsed.Execute(outFile, r.config); err != nil {
			return errors.Join(err)
		}
		return nil
	})
	if err != nil {
		return errors.Join(err)
	}
	return nil
}
