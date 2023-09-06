package ecs

import (
	"context"
	"errors"
	"fmt"
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

func NewProvider(ctx context.Context, options *options.Options, logs log.Logger) (*EcsProvider, error) {
	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	// create provider
	provider := &EcsProvider{
		Config:    options,
		AwsConfig: cfg,
		Log:       logs,

		client: ecs.NewFromConfig(cfg),
	}

	return provider, nil
}

type EcsProvider struct {
	Config    *options.Options
	AwsConfig aws.Config
	Log       log.Logger

	client *ecs.Client
}

func (p *EcsProvider) TargetArchitecture(ctx context.Context, workspaceId string) (string, error) {
	return p.Config.ClusterArchitecture, nil
}

func (p *EcsProvider) StartTask(ctx context.Context, workspaceId string) error {
	// noop operation if running on fargate
	if p.Config.LaunchType == string(types.LaunchTypeFargate) {
		return nil
	}

	return p.startTask(ctx, workspaceId)
}

func (p *EcsProvider) StopTask(ctx context.Context, workspaceId string) error {
	// noop operation if running on fargate
	if p.Config.LaunchType == string(types.LaunchTypeFargate) {
		return nil
	}
	
	// stop the task
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

	return nil
}

func (p *EcsProvider) RunTask(ctx context.Context, workspaceId string, runOptions *driver.RunOptions) error {
	err := p.registerTaskDefinition(ctx, workspaceId, runOptions)
	if err != nil {
		return err
	}

	err = p.startTask(ctx, workspaceId)
	if err != nil {
		_ = p.deleteTaskDefinition(ctx, workspaceId)
		return err
	}

	return nil
}

func (p *EcsProvider) FindTask(ctx context.Context, workspaceId string) (*config.ContainerDetails, error) {
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

func (p *EcsProvider) DeleteTask(ctx context.Context, workspaceId string) error {
	// stop the task
	err := p.StopTask(ctx, workspaceId)
	if err != nil {
		return err
	}

	// TODO: delete ecs volume?

	// delete task definition
	err = p.deleteTaskDefinition(ctx, workspaceId)
	if err != nil {
		return err
	}

	return nil
}

func (p *EcsProvider) getTaskID(ctx context.Context, workspaceId string) (*types.Task, error) {
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

func (p *EcsProvider) startTask(ctx context.Context, workspaceId string) error {
	taskDefinitionID, err := p.getTaskDefinitionArn(ctx, workspaceId)
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
		LaunchType:           types.LaunchType(p.Config.LaunchType),
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        []string{p.Config.SubnetID},
				SecurityGroups: securityGroups,
				AssignPublicIp: types.AssignPublicIp(p.Config.AssignPublicIp),
			},
		},
		Tags: getTags(workspaceId),
	})
	if err != nil {
		return fmt.Errorf("run task: %w", err)
	} else if len(taskOutput.Failures) > 0 {
		return fmt.Errorf("run task failure: %w", errors.New(*taskOutput.Failures[0].Reason))
	}

	return nil
}
