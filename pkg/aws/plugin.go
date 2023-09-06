package aws

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/session-manager-plugin/src/datachannel"
	"github.com/aws/session-manager-plugin/src/log"
	"github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session"
	_ "github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session/portsession"
	_ "github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session/shellsession"
	"github.com/google/uuid"
	"github.com/loft-sh/devpod-provider-ecs/pkg/options"
)

func (p *AwsProvider) StartSession(target string, port int) error {
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
