package options

import (
	"fmt"
	"os"
)

var DefaultSSHPort int = 19583

type Options struct {
	DevContainerID string

	ClusterID           string
	ClusterArchitecture string

	SubnetID        string
	SecurityGroupID string

	TaskRoleARN      string
	ExecutionRoleARN string

	TaskCpu    string
	TaskMemory string

	LaunchType     string
	AssignPublicIp string
}

func FromEnv() (*Options, error) {
	retOptions := &Options{}

	var err error

	// required
	retOptions.DevContainerID, err = fromEnvOrError("DEVCONTAINER_ID")
	if err != nil {
		return nil, err
	}
	retOptions.ClusterID, err = fromEnvOrError("CLUSTER_ID")
	if err != nil {
		return nil, err
	}
	retOptions.SubnetID, err = fromEnvOrError("SUBNET_ID")
	if err != nil {
		return nil, err
	}
	retOptions.TaskRoleARN, err = fromEnvOrError("TASK_ROLE_ARN")
	if err != nil {
		return nil, err
	}
	retOptions.ExecutionRoleARN, err = fromEnvOrError("EXECUTION_ROLE_ARN")
	if err != nil {
		return nil, err
	}
	retOptions.ClusterArchitecture, err = fromEnvOrError("CLUSTER_ARCHITECTURE")
	if err != nil {
		return nil, err
	}
	retOptions.TaskCpu, err = fromEnvOrError("TASK_CPU")
	if err != nil {
		return nil, err
	}
	retOptions.TaskMemory, err = fromEnvOrError("TASK_MEMORY")
	if err != nil {
		return nil, err
	}
	retOptions.LaunchType, err = fromEnvOrError("LAUNCH_TYPE")
	if err != nil {
		return nil, err
	}
	retOptions.AssignPublicIp, err = fromEnvOrError("ASSIGN_PUBLIC_IP")
	if err != nil {
		return nil, err
	}

	// optional
	retOptions.SecurityGroupID = os.Getenv("SECURITY_GROUP_ID")

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
