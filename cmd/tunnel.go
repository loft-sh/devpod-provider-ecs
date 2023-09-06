package cmd

import (
	"context"

	"github.com/loft-sh/devpod-provider-ecs/pkg/ecs"
	"github.com/loft-sh/devpod-provider-ecs/pkg/options"
	"github.com/loft-sh/log"
	"github.com/spf13/cobra"
)

type TunnelCmd struct {
	Target string

	Port int
}

// NewTunnelCmd returns a new command
func NewTunnelCmd() *cobra.Command {
	cmd := &TunnelCmd{}
	cobraCmd := &cobra.Command{
		Use:           "tunnel",
		Short:         "Creates a ssh tunnel through aws ssm to the ecs container",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.Run(cobraCmd.Context())
		},
	}

	cobraCmd.Flags().StringVar(&cmd.Target, "target", "", "The target to connect to.")
	cobraCmd.Flags().IntVar(&cmd.Port, "port", options.DefaultSSHPort, "The port to use where the ssh server is running")
	return cobraCmd
}

func (cmd *TunnelCmd) Run(ctx context.Context) error {
	awsOptions, err := options.FromEnv()
	if err != nil {
		return err
	}

	ecsProvider, err := ecs.NewProvider(ctx, awsOptions, log.Default.ErrorStreamOnly())
	if err != nil {
		return err
	}

	return ecsProvider.StartSession(cmd.Target, cmd.Port)
}
