package options

import (
	"fmt"
	"os"
)

var DefaultSSHPort int = 19583

type Options struct {
	ClusterID string

	SubnetID        string
	SecurityGroupID string

	ExecutionRoleARN string
}

func FromEnv(init bool) (*Options, error) {
	retOptions := &Options{}

	return retOptions, nil
}

func fromEnvOrError(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf(
			"couldn't find option %s in environment, please make sure %s is defined",
			name,
			name,
		)
	}

	return val, nil
}
