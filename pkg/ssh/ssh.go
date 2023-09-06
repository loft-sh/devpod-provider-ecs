package ssh

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/loft-sh/devpod/pkg/ssh"
	"github.com/loft-sh/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func ExecuteCommand(ctx context.Context, command string, stdin io.Reader, stdout, stderr io.Writer, log log.Logger) error {
	// create context
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	stdinReader, stdinWriter, err := os.Pipe()
	if err != nil {
		return err
	}
	defer stdinWriter.Close()

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		return err
	}
	defer stdoutWriter.Close()

	tunnelChan := make(chan error, 1)
	go func() {
		writer := log.ErrorStreamOnly().Writer(logrus.InfoLevel, false)
		defer writer.Close()

		tunnelChan <- startProxyCommand(cancelCtx, stdinReader, stdoutWriter, writer)
	}()

	// connect to container
	containerChan := make(chan error, 1)
	go func() {
		// start ssh client as root / default user
		sshClient, err := ssh.StdioClient(stdoutReader, stdinWriter, false)
		if err != nil {
			containerChan <- errors.Wrap(err, "create ssh client")
			return
		}

		defer sshClient.Close()
		defer cancel()

		containerChan <- ssh.Run(cancelCtx, sshClient, command, stdin, stdout, stderr)
	}()

	// wait for result
	select {
	case err := <-containerChan:
		return errors.Wrap(err, "ssh into container")
	case err := <-tunnelChan:
		return errors.Wrap(err, "connect to ssm")
	}
}

func startProxyCommand(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(
		ctx,
		executable,
		"tunnel",
		"--target", "ecs:fabian-cluster_8e599549198849e3b9eb645a56572b24_8e599549198849e3b9eb645a56572b24-504127535",
	)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}
