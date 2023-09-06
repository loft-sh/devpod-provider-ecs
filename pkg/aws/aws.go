package aws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/loft-sh/devpod-provider-ecs/pkg/options"
	"github.com/loft-sh/devpod/pkg/devcontainer/config"
	"github.com/loft-sh/devpod/pkg/driver"
	"github.com/loft-sh/log"
)

func NewProvider(ctx context.Context, logs log.Logger) (*AwsProvider, error) {
	config, err := options.FromEnv(false)
	if err != nil {
		return nil, err
	}

	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	// create provider
	provider := &AwsProvider{
		Config:    config,
		AwsConfig: cfg,
		Log:       logs,

		client: ecs.NewFromConfig(cfg),
	}

	return provider, nil
}

type AwsProvider struct {
	Config    *options.Options
	AwsConfig aws.Config
	Log       log.Logger

	client *ecs.Client
}

func (p *AwsProvider) ExecuteCommand(ctx context.Context, workspaceId, command string, stdout io.Writer, stdin io.Reader) error {
	task, err := p.getTaskID(ctx, workspaceId)
	if err != nil {
		return err
	} else if task == nil {
		return fmt.Errorf("no task for workspace %s found", workspaceId)
	}

	taskArnSplitted := strings.Split(*task.TaskArn, "/")
	fmt.Println("ecs:" + p.Config.ClusterID + "_" + taskArnSplitted[len(taskArnSplitted)-1] + "_" + *task.Containers[0].RuntimeId)
	return nil
}

func (p *AwsProvider) FindTask(ctx context.Context, workspaceId string) (*config.ContainerDetails, error) {
	task, err := p.getTaskID(ctx, workspaceId)
	if err != nil {
		return nil, err
	} else if task == nil {
		return nil, nil
	}

	// get labels
	taskDefinition, err := p.client.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: task.TaskDefinitionArn,
	})
	if err != nil {
		return nil, fmt.Errorf("describe task definition: %w", err)
	}
	labels := taskDefinition.TaskDefinition.ContainerDefinitions[0].DockerLabels

	// status
	status := "created"
	if task.LastStatus != nil && strings.ToUpper(*task.LastStatus) == string(types.DesiredStatusRunning) {
		status = "running"
	} else if task.LastStatus != nil && strings.ToUpper(*task.LastStatus) == string(types.DesiredStatusStopped) {
		status = "exited"
	}

	// started at
	startedAt := ""
	if task.StartedAt != nil {
		startedAt = task.StartedAt.String()
	}

	return &config.ContainerDetails{
		ID:      *task.TaskArn,
		Created: task.CreatedAt.String(),
		State: config.ContainerDetailsState{
			Status:    status,
			StartedAt: startedAt,
		},
		Config: config.ContainerDetailsConfig{
			Labels: labels,
		},
	}, nil
}

func (p *AwsProvider) DeleteTask(ctx context.Context, workspaceId string) error {
	task, err := p.getTaskID(ctx, workspaceId)
	if err != nil {
		return err
	} else if task != nil {
		// delete the task
		_, err = p.client.StopTask(ctx, &ecs.StopTaskInput{
			Task:    task.TaskArn,
			Cluster: options.Ptr(p.Config.ClusterID),
		})
		if err != nil {
			return fmt.Errorf("stop task: %w", err)
		}
	}

	// delete task definition
	err = p.deleteTaskDefinition(ctx, workspaceId)
	if err != nil {
		return err
	}

	return nil
}

func (p *AwsProvider) getTaskID(ctx context.Context, workspaceId string) (*types.Task, error) {
	runningTaskArns, err := p.client.ListTasks(ctx, &ecs.ListTasksInput{
		Cluster:       options.Ptr(p.Config.ClusterID),
		Family:        options.Ptr("devpod-" + workspaceId),
		DesiredStatus: types.DesiredStatusRunning,
		MaxResults:    options.Ptr(int32(10)),
	})
	if err != nil {
		return nil, fmt.Errorf("list running tasks: %w", err)
	}

	stoppedTaskArns, err := p.client.ListTasks(ctx, &ecs.ListTasksInput{
		Cluster:       options.Ptr(p.Config.ClusterID),
		Family:        options.Ptr("devpod-" + workspaceId),
		DesiredStatus: types.DesiredStatusStopped,
		MaxResults:    options.Ptr(int32(10)),
	})
	if err != nil {
		return nil, fmt.Errorf("list stopped tasks: %w", err)
	}

	taskArns := append(runningTaskArns.TaskArns, stoppedTaskArns.TaskArns...)
	if len(taskArns) == 0 {
		return nil, nil
	}

	tasks, err := p.client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Tasks:   taskArns,
		Cluster: options.Ptr(p.Config.ClusterID),
	})
	if err != nil {
		return nil, fmt.Errorf("describe tasks: %w", err)
	} else if len(tasks.Failures) > 0 {
		return nil, fmt.Errorf("describe tasks failures: %s", *tasks.Failures[0].Reason)
	} else if len(tasks.Tasks) == 0 {
		return nil, nil
	}

	// sort tasks by revision
	sort.SliceStable(tasks.Tasks, func(i, j int) bool {
		return tasks.Tasks[i].CreatedAt.Unix() > tasks.Tasks[j].CreatedAt.Unix()
	})

	return &tasks.Tasks[0], nil
}

func (p *AwsProvider) RunTask(ctx context.Context, workspaceId string, runOptions *driver.RunOptions) error {
	taskDefinitionID, err := p.registerTaskDefinition(ctx, workspaceId, runOptions)
	if err != nil {
		return err
	}

	securityGroups := []string{}
	if p.Config.SecurityGroupID != "" {
		securityGroups = append(securityGroups, p.Config.SecurityGroupID)
	}

	taskOutput, err := p.client.RunTask(ctx, &ecs.RunTaskInput{
		TaskDefinition:       options.Ptr(taskDefinitionID),
		Cluster:              options.Ptr(p.Config.ClusterID),
		Count:                options.Ptr(int32(1)),
		EnableExecuteCommand: true,
		LaunchType:           types.LaunchTypeFargate,
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        []string{p.Config.SubnetID},
				SecurityGroups: securityGroups,
				AssignPublicIp: types.AssignPublicIpEnabled,
			},
		},
		Tags: getTags(workspaceId),
	})
	if err != nil {
		_ = p.deleteTaskDefinition(ctx, workspaceId)
		return fmt.Errorf("run task: %w", err)
	} else if len(taskOutput.Failures) > 0 {
		_ = p.deleteTaskDefinition(ctx, workspaceId)
		return fmt.Errorf("run task failure: %w", errors.New(*taskOutput.Failures[0].Reason))
	}

	return nil
}

func (p *AwsProvider) registerTaskDefinition(ctx context.Context, workspaceId string, runOptions *driver.RunOptions) (string, error) {
	taskDefinitionID := "devpod-" + workspaceId

	// delete existing task definition
	err := p.deleteTaskDefinition(ctx, workspaceId)
	if err != nil {
		return "", fmt.Errorf("delete existing task definition: %w", err)
	}

	// get container definition
	containerDefinition, err := getContainerDefinition(runOptions)
	if err != nil {
		return "", fmt.Errorf("get container definition: %w", err)
	}

	// register task definition
	_, err = p.client.RegisterTaskDefinition(ctx, &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []types.ContainerDefinition{
			containerDefinition,
		},
		TaskRoleArn:      options.Ptr(p.Config.ExecutionRoleARN),
		ExecutionRoleArn: options.Ptr(p.Config.ExecutionRoleARN),
		Family:           options.Ptr(taskDefinitionID),
		Cpu:              options.Ptr(".5 vcpu"),
		Memory:           options.Ptr("1 gb"),
		NetworkMode:      types.NetworkModeAwsvpc,
		RequiresCompatibilities: []types.Compatibility{
			types.CompatibilityFargate,
		},
		Tags: getTags(workspaceId),
	})
	if err != nil {
		return "", err
	}

	// retrieve the task definition arn
	return p.getTaskDefinitionArn(ctx, workspaceId)
}

func (p *AwsProvider) deleteTaskDefinition(ctx context.Context, workspaceId string) error {
	taskDefinitionID := "devpod-" + workspaceId

	// list existing task definitions
	output, err := p.client.ListTaskDefinitions(ctx, &ecs.ListTaskDefinitionsInput{
		FamilyPrefix: options.Ptr(taskDefinitionID),
		MaxResults:   options.Ptr(int32(10)),
	})
	if err != nil {
		return err
	} else if len(output.TaskDefinitionArns) > 0 {
		// deregister task definitions
		for _, taskDefinition := range output.TaskDefinitionArns {
			_, err = p.client.DeregisterTaskDefinition(ctx, &ecs.DeregisterTaskDefinitionInput{
				TaskDefinition: &taskDefinition,
			})
			if err != nil {
				return fmt.Errorf("deregister task definition %s: %w", taskDefinition, err)
			}
		}

		// delete existing task definitions
		output, err := p.client.DeleteTaskDefinitions(ctx, &ecs.DeleteTaskDefinitionsInput{
			TaskDefinitions: output.TaskDefinitionArns,
		})
		if err != nil {
			return err
		} else if len(output.Failures) > 0 {
			return errors.New(*output.Failures[0].Reason)
		}
	}

	return nil
}

func (p *AwsProvider) getTaskDefinitionArn(ctx context.Context, workspaceId string) (string, error) {
	taskDefinitionID := "devpod-" + workspaceId

	// list existing task definitions
	output, err := p.client.ListTaskDefinitions(ctx, &ecs.ListTaskDefinitionsInput{
		FamilyPrefix: options.Ptr(taskDefinitionID),
		MaxResults:   options.Ptr(int32(10)),
	})
	if err != nil {
		return "", err
	} else if len(output.TaskDefinitionArns) != 1 {
		return "", fmt.Errorf("unexpected amount of task definitions: %d, expected 1", len(output.TaskDefinitionArns))
	}

	return output.TaskDefinitionArns[0], nil
}

func getContainerDefinition(runOptions *driver.RunOptions) (types.ContainerDefinition, error) {
	retDefinition := types.ContainerDefinition{
		Name:      options.Ptr("devpod"),
		Image:     &runOptions.Image,
		Essential: options.Ptr(true),
		LinuxParameters: &types.LinuxParameters{
			InitProcessEnabled: options.Ptr(true),
		},
	}
	if len(runOptions.Labels) > 0 {
		retDefinition.DockerLabels = config.ListToObject(runOptions.Labels)
	}
	if len(runOptions.Env) > 0 {
		for k, v := range runOptions.Env {
			retDefinition.Environment = append(retDefinition.Environment, types.KeyValuePair{
				Name:  options.Ptr(k),
				Value: options.Ptr(v),
			})
		}
	}
	retDefinition.EntryPoint = []string{runOptions.Entrypoint}
	retDefinition.Command = runOptions.Cmd
	if runOptions.User != "" {
		retDefinition.User = &runOptions.User
	}
	retDefinition.DockerSecurityOptions = runOptions.SecurityOpt
	retDefinition.Privileged = runOptions.Privileged

	// TODO: volumes

	return retDefinition, nil
}

func getTags(workspaceId string) []types.Tag {
	return []types.Tag{
		{
			Key:   options.Ptr("devpod-workspace-id"),
			Value: options.Ptr(workspaceId),
		},
	}
}
