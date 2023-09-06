package cmd

import (
	"context"
	"os"

	"github.com/loft-sh/devpod-provider-ecs/pkg/ecs"
	"github.com/loft-sh/devpod-provider-ecs/pkg/options"
	"github.com/loft-sh/log"
	"github.com/spf13/cobra"
)

// CommandCmd holds the cmd flags
type CommandCmd struct{}

// NewCommandCmd defines a command
func NewCommandCmd() *cobra.Command {
	cmd := &CommandCmd{}
	commandCmd := &cobra.Command{
		Use:   "command",
		Short: "Command a container",
		RunE: func(_ *cobra.Command, args []string) error {
			options, err := options.FromEnv()
			if err != nil {
				return err
			}

			return cmd.Run(context.Background(), options, log.Default.ErrorStreamOnly())
		},
	}

	return commandCmd
}

// Run runs the command logic
func (cmd *CommandCmd) Run(ctx context.Context, options *options.Options, log log.Logger) error {
	ecsProvider, err := ecs.NewProvider(ctx, options, log)
	if err != nil {
		return err
	}

	return ecsProvider.ExecuteCommand(
		ctx,
		options.DevContainerID,
		os.Getenv("DEVCONTAINER_USER"),
		os.Getenv("DEVCONTAINER_COMMAND"),
		os.Stdin,
		os.Stdout,
		os.Stderr,
	)
}
