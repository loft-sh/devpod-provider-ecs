package cmd

import (
	"context"
	"os"

	"github.com/loft-sh/devpod-provider-ecs/pkg/aws"
	"github.com/loft-sh/devpod-provider-ecs/pkg/ssh"
	"github.com/loft-sh/log"
	"github.com/spf13/cobra"
)

// NewTestCmd returns a new command
func NewTestCmd() *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           "test",
		Short:         "TODO: remove this",
		Hidden:        true,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			provider, err := aws.NewProvider(context.Background(), log.Default)
			if err != nil {
				panic(err)
			}

			ctx := context.Background()
			//workspaceId := "fabian-test-workspace-id"
			provider.Config.ClusterID = "fabian-cluster"
			provider.Config.SubnetID = "subnet-9c749bd0"
			provider.Config.ExecutionRoleARN = "arn:aws:iam::977114253874:role/ecsTaskExecutionRole"

			//err = provider.DeleteTask(ctx, workspaceId)
			err = ssh.ExecuteCommand(ctx, "ls /", os.Stdin, os.Stdout, os.Stderr, log.Default)
			/*err = provider.ExecuteCommand(
				ctx,
				workspaceId,
				"sh -c 'ls'",
				os.Stdout,
				strings.NewReader("pongadfhjdfshkjsfdhk\n"),
			)*/
			/*err = provider.RunTask(ctx, workspaceId, &driver.RunOptions{
				Image: "ghcr.io/loft-sh/dockerless:0.1.2",
				User:  "root",
				Env: map[string]string{
					"FABI": "test123",
				},
			})*/
			if err != nil {
				panic(err)
			}
			/*
				err = provider.RunTask(context.Background(), workspaceId, &driver.RunOptions{
					Image:      "alpine",
					User:       "root",
					Entrypoint: "/bin/sh",
					Cmd:        []string{"-c", "sleep 100000"},
					Env: map[string]string{
						"FABI": "test123",
					},
				})
				if err != nil {
					panic(err)
				}*/

			return nil
		},
	}

	return cobraCmd
}
