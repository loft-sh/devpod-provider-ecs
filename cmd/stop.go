package cmd

import (
	"context"

	"github.com/loft-sh/devpod-provider-ecs/pkg/ecs"
	"github.com/loft-sh/devpod-provider-ecs/pkg/options"
	"github.com/loft-sh/log"
	"github.com/spf13/cobra"
)

// StopCmd holds the cmd flags
type StopCmd struct{}

// NewStopCmd defines a command
func NewStopCmd() *cobra.Command {
	cmd := &StopCmd{}
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a container",
		RunE: func(_ *cobra.Command, args []string) error {
			options, err := options.FromEnv()
			if err != nil {
				return err
			}

			return cmd.Run(context.Background(), options, log.Default)
		},
	}

	return stopCmd
}

// Run runs the command logic
func (cmd *StopCmd) Run(ctx context.Context, options *options.Options, log log.Logger) error {
	ecsProvider, err := ecs.NewProvider(ctx, options, log)
	if err != nil {
		return err
	}

	return ecsProvider.StopTask(ctx, options.DevContainerID)
}
