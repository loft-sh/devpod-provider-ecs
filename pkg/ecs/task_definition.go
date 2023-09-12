package ecs

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/loft-sh/devpod-provider-ecs/pkg/hash"
	"github.com/loft-sh/devpod-provider-ecs/pkg/inject"
	"github.com/loft-sh/devpod-provider-ecs/pkg/options"
	"github.com/loft-sh/devpod/pkg/devcontainer/config"
	"github.com/loft-sh/devpod/pkg/driver"
)

func (p *EcsProvider) registerTaskDefinition(ctx context.Context, workspaceId string, runOptions *driver.RunOptions) error {
	taskDefinitionID := "devpod-" + workspaceId

	// delete existing task definition
	err := p.deleteTaskDefinition(ctx, workspaceId)
	if err != nil {
		return fmt.Errorf("delete existing task definition: %w", err)
	}

	// get container definition
	containerDefinition, err := p.getContainerDefinition(workspaceId, runOptions)
	if err != nil {
		return fmt.Errorf("get container definition: %w", err)
	}

	// make sure we have a value for the role arn
	if p.Config.TaskRoleARN == "" || p.Config.ExecutionRoleARN == "" {
		roleArn, err := p.createIamRole(ctx)
		if err != nil {
			return err
		}

		if p.Config.TaskRoleARN == "" {
			p.Config.TaskRoleARN = roleArn
		}
		if p.Config.ExecutionRoleARN == "" {
			p.Config.ExecutionRoleARN = roleArn
		}
	}

	// create task definition
	taskDefinition := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []types.ContainerDefinition{
			containerDefinition,
		},
		TaskRoleArn:      options.Ptr(p.Config.TaskRoleARN),
		ExecutionRoleArn: options.Ptr(p.Config.ExecutionRoleARN),
		Family:           options.Ptr(taskDefinitionID),
		Cpu:              options.Ptr(p.Config.TaskCpu),
		Memory:           options.Ptr(p.Config.TaskMemory),
		NetworkMode:      types.NetworkModeAwsvpc,
		RequiresCompatibilities: []types.Compatibility{
			types.Compatibility(p.Config.LaunchType),
		},
		Tags: getTags(workspaceId),
	}

	// add volumes
	if p.Config.LaunchType != string(types.LaunchTypeFargate) {
		dockerVolumeConfiguration := &types.DockerVolumeConfiguration{
			Autoprovision: options.Ptr(true),
			Driver:        options.Ptr("local"),
			Scope:         "shared",
		}
		taskDefinition.Volumes = append(taskDefinition.Volumes, types.Volume{
			Name:                      options.Ptr("devpod-" + workspaceId),
			DockerVolumeConfiguration: dockerVolumeConfiguration,
		})
		for _, mount := range runOptions.Mounts {
			if mount.Source == "" || mount.Target == "" {
				continue
			}

			taskDefinition.Volumes = append(taskDefinition.Volumes, types.Volume{
				Name:                      options.Ptr(volumeName(workspaceId, mount.Source)),
				DockerVolumeConfiguration: dockerVolumeConfiguration,
			})
		}
	}

	// register task definition
	_, err = p.client.RegisterTaskDefinition(ctx, taskDefinition)
	if err != nil {
		return err
	}

	return nil
}

func (p *EcsProvider) getTaskDefinitionArn(ctx context.Context, workspaceId string) (string, error) {
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

func (p *EcsProvider) getContainerDefinition(workspaceId string, runOptions *driver.RunOptions) (types.ContainerDefinition, error) {
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

	entrypoint, cmd, err := inject.GetContainerEntrypoint([]string{runOptions.Entrypoint}, runOptions.Cmd)
	if err != nil {
		return types.ContainerDefinition{}, err
	}

	retDefinition.EntryPoint = entrypoint
	retDefinition.Command = cmd
	if runOptions.User != "" {
		retDefinition.User = &runOptions.User
	}
	retDefinition.DockerSecurityOptions = runOptions.SecurityOpt
	retDefinition.Privileged = runOptions.Privileged

	// mount points
	if p.Config.LaunchType != string(types.LaunchTypeFargate) {
		retDefinition.MountPoints = append(retDefinition.MountPoints, types.MountPoint{
			ContainerPath: options.Ptr("/workspaces"),
			SourceVolume:  options.Ptr("devpod-" + workspaceId),
		})
		for _, mount := range runOptions.Mounts {
			if mount.Source == "" || mount.Target == "" {
				continue
			}

			retDefinition.MountPoints = append(retDefinition.MountPoints, types.MountPoint{
				ContainerPath: options.Ptr(mount.Target),
				SourceVolume:  options.Ptr(volumeName(workspaceId, mount.Source)),
			})
		}
	}

	return retDefinition, nil
}

func volumeName(workspaceId, source string) string {
	return "devpod-" + workspaceId + "-" + hash.String(source)[:5]
}

func (p *EcsProvider) deleteTaskDefinition(ctx context.Context, workspaceId string) error {
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
		p.Log.Info("Deleting task definition...")
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

func getTags(workspaceId string) []types.Tag {
	return []types.Tag{
		{
			Key:   options.Ptr("devpod-workspace-id"),
			Value: options.Ptr(workspaceId),
		},
	}
}
