package cmd

import (
	"fmt"

	"github.com/loft-sh/devpod-provider-ecs/pkg/options"
	sshserver "github.com/loft-sh/devpod/pkg/ssh/server"
	"github.com/loft-sh/log"
	"github.com/spf13/cobra"
)

type EntrypointCmd struct {
	Entrypoint []string
	Cmd        []string

	User string

	Port int
}

// NewEntrypointCmd returns a new start command
func NewEntrypointCmd() *cobra.Command {
	cmd := &EntrypointCmd{}
	cobraCmd := &cobra.Command{
		Use:           "entrypoint",
		Short:         "Starts the container with an ssh server in the background",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.Run()
		},
	}

	cobraCmd.Flags().StringVar(&cmd.User, "user", "", "The container user to run the entrypoint with.")
	cobraCmd.Flags().StringArrayVar(&cmd.Entrypoint, "entrypoint", []string{}, "The entrypoint to use.")
	cobraCmd.Flags().StringArrayVar(&cmd.Cmd, "cmd", []string{}, "The cmds to use.")
	cobraCmd.Flags().IntVar(&cmd.Port, "port", options.DefaultSSHPort, "The default port to use for the ssh server")
	return cobraCmd
}

func (cmd *EntrypointCmd) Run() error {
	address := fmt.Sprintf("127.0.0.1:%d", cmd.Port)
	server, err := sshserver.NewServer(address, nil, nil, log.Default)
	if err != nil {
		return err
	}

	log.Default.Infof("Listen and serve on: %s", address)
	err = server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
