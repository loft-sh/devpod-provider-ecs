package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/loft-sh/devpod-provider-ecs/pkg/options"
	sshserver "github.com/loft-sh/devpod/pkg/ssh/server"
	"github.com/loft-sh/log"
	"github.com/spf13/cobra"
)

type EntrypointCmd struct {
	Entrypoint string
	Cmd        string

	Port int
}

// NewEntrypointCmd returns a new command
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

	cobraCmd.Flags().StringVar(&cmd.Entrypoint, "entrypoint", "", "Base64 encoded json string with an entrypoint to execute")
	cobraCmd.Flags().StringVar(&cmd.Cmd, "cmd", "", "Base64 encoded json string with cmd to execute")
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
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			log.Default.Fatal("SSH server failed: %v", err)
		} else {
			log.Default.Fatal("SSH server ended unexpectedly")
		}
	}()

	args := []string{}
	if cmd.Entrypoint != "" {
		entrypoint, err := decodeStrArray(cmd.Entrypoint)
		if err != nil {
			return fmt.Errorf("decode entrypoint: %w", err)
		}

		args = append(args, entrypoint...)
	}
	if cmd.Cmd != "" {
		cmd, err := decodeStrArray(cmd.Cmd)
		if err != nil {
			return fmt.Errorf("decode cmd: %w", err)
		}

		args = append(args, cmd...)
	}

	// run entrypoint?
	if len(args) == 0 {
		// wait indefinitely
		select {}
		return nil
	}

	// run entrypoint
	entrypointCmd := exec.Command(args[0], args[1:]...)
	entrypointCmd.Stdout = os.Stdout
	entrypointCmd.Stdin = os.Stdin
	entrypointCmd.Stderr = os.Stderr
	err = entrypointCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func decodeStrArray(payload string) ([]string, error) {
	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	strArr := []string{}
	err = json.Unmarshal(decoded, &strArr)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", string(decoded), err)
	}

	return strArr, nil
}
