package inject

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/loft-sh/devpod-provider-ecs/pkg/version"
)

//go:embed inject.sh
var Script string

const BaseURL = "https://github.com/loft-sh/devpod-provider-ecs/releases/download/%s/devpod-provider-ecs-linux-%s"

const LatestBaseURL = "https://github.com/loft-sh/devpod-provider-ecs/releases/latest/download/devpod-provider-ecs-linux-%s"

func GetContainerEntrypoint(entrypoint []string, cmd []string) ([]string, []string, error) {
	downloadAmd := ""
	downloadArm := ""
	if version.Version == "latest" {
		downloadAmd = fmt.Sprintf(LatestBaseURL, "amd64")
		downloadArm = fmt.Sprintf(LatestBaseURL, "arm64")
	} else {
		downloadAmd = fmt.Sprintf(BaseURL, version.Version, "amd64")
		downloadArm = fmt.Sprintf(BaseURL, version.Version, "arm64")
	}

	injectScript, err := FillTemplate(Script, map[string]string{
		"DownloadAmd":     downloadAmd,
		"DownloadArm":     downloadArm,
		"InstallFilename": "devpod-provider-ecs",
		"InstallDir":      "/workspaces",
		"Command":         "/workspaces/devpod-provider-ecs entrypoint",
	})
	if err != nil {
		return nil, nil, err
	}

	return []string{"sh"}, []string{"-c", injectScript}, nil
}

func FillTemplate(templateString string, vars interface{}) (string, error) {
	t, err := template.New("gotmpl").Parse(templateString)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	err = t.Execute(&buf, vars)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
