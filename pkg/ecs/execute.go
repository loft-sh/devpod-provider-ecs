package ecs

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/session-manager-plugin/src/datachannel"
	"github.com/aws/session-manager-plugin/src/log"
	"github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session"
	_ "github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session/portsession"
	_ "github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session/shellsession"
	"github.com/google/uuid"
	"github.com/loft-sh/devpod-provider-ecs/pkg/options"
	"github.com/loft-sh/devpod/pkg/ssh"
	loftlog "github.com/loft-sh/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (p *EcsProvider) ExecuteCommand(ctx context.Context, workspaceId, user, command string, stdin io.Reader, stdout, stderr io.Writer) error {
	task, err := p.getTaskID(ctx, workspaceId)
	if err != nil {
		return err
	} else if task == nil {
		return fmt.Errorf("no task for workspace %s found", workspaceId)
	}

	target := "ecs:" + getIDFromArn(p.Config.ClusterID) + "_" + getIDFromArn(*task.TaskArn) + "_" + *task.Containers[0].RuntimeId
	return executeCommand(ctx, target, user, command, stdin, stdout, stderr, p.Log)
}

func executeCommand(ctx context.Context, target, user, command string, stdin io.Reader, stdout, stderr io.Writer, log loftlog.Logger) error {
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

		tunnelChan <- startProxyCommand(cancelCtx, target, stdinReader, stdoutWriter, writer)
	}()

	// connect to container
	containerChan := make(chan error, 1)
	go func() {
		// start ssh client as root / default user
		sshClient, err := ssh.StdioClientWithUser(stdoutReader, stdinWriter, user, false)
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

func startProxyCommand(ctx context.Context, target string, stdin io.Reader, stdout, stderr io.Writer) error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(
		ctx,
		executable,
		"tunnel",
		"--target", target,
	)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func getIDFromArn(arn string) string {
	if !strings.HasPrefix(arn, "arn:") {
		return arn
	}

	taskArnSplitted := strings.Split(arn, "/")
	return taskArnSplitted[len(taskArnSplitted)-1]
}

func (p *EcsProvider) StartSession(target string, port int) error {
	out, err := ssm.NewFromConfig(p.AwsConfig).StartSession(context.Background(), &ssm.StartSessionInput{
		Target:       options.Ptr(target),
		DocumentName: options.Ptr("AWS-StartSSHSession"),
		Parameters: map[string][]string{
			"portNumber": {strconv.Itoa(port)},
		},
	})
	if err != nil {
		return err
	}

	ssmSession := new(session.Session)
	ssmSession.SessionId = *out.SessionId
	ssmSession.StreamUrl = *out.StreamUrl
	ssmSession.TokenValue = *out.TokenValue
	ssmSession.ClientId = uuid.NewString()
	ssmSession.TargetId = target
	ssmSession.DataChannel = &datachannel.DataChannel{}
	return ssmSession.Execute(log.Logger(false, ssmSession.ClientId))
}
