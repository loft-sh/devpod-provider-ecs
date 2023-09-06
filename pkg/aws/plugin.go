package aws

/*
import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/session-manager-plugin/src/datachannel"
	"github.com/aws/session-manager-plugin/src/log"
	"github.com/aws/session-manager-plugin/src/message"
	"github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session"
	_ "github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session/portsession"
	_ "github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session/shellsession"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func PluginSession(out *types.Session, writer io.Writer, reader io.Reader) error {
	ssmSession := &session.Session{}
	ssmSession.SessionId = *out.SessionId
	ssmSession.StreamUrl = *out.StreamUrl
	ssmSession.TokenValue = *out.TokenValue
	ssmSession.ClientId = uuid.NewString()
	ssmSession.DataChannel = &datachannel.DataChannel{}
	ssmSession.StopChan = make(chan struct{})

	logger := log.Logger(false, ssmSession.ClientId)
	err := ssmSession.Execute(logger, writer)
	if err != nil {
		return err
	}

	// Send compressed archive to client
	go func() {
		defer fmt.Println("DONE STDIN")

		err = pipeStdin(reader, ssmSession, logger)
		if err != nil {
			fmt.Println(err)
		}
	}()

	// wait for command to finish
	<-ssmSession.StopChan
	return nil
}

func pipeStdin(reader io.Reader, ssmSession *session.Session, logger log.T) error {
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			err := ssmSession.DataChannel.SendInputDataMessage(logger, message.Output, buf[:n])
			if err != nil {
				return errors.Wrap(err, "stream send")
			}
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return errors.Wrap(err, "read file")
		}
	}

	return nil
}
*/
