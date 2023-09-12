package ecs

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/loft-sh/devpod-provider-ecs/pkg/options"
)

var (
	devPodRoleName   = "devpod-ecs-role"
	devPodPolicyName = "devpod-ecs-policy"
)

func (p *EcsProvider) createIamRole(ctx context.Context) (string, error) {
	// check for role
	iamClient := iam.NewFromConfig(p.AwsConfig)
	role, err := iamClient.GetRole(ctx, &iam.GetRoleInput{
		RoleName: &devPodRoleName,
	})
	if err != nil {
		var re *awshttp.ResponseError
		if !errors.As(err, &re) || re.HTTPStatusCode() != http.StatusNotFound {
			return "", err
		}
	} else {
		return *role.Role.Arn, nil
	}

	// create policy
	p.Log.Infof("Create iam policy %s...", devPodPolicyName)
	policyOutput, err := iamClient.CreatePolicy(ctx, &iam.CreatePolicyInput{
		PolicyName: &devPodPolicyName,
		PolicyDocument: options.Ptr(`{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "ecs:ExecuteCommand",
                "ssmmessages:CreateControlChannel",
                "ssmmessages:CreateDataChannel",
                "ssmmessages:OpenControlChannel",
                "ssmmessages:OpenDataChannel",
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Effect": "Allow",
            "Resource": "*"
        }
    ]
}`),
	})
	if err != nil {
		return "", fmt.Errorf("create policy: %w", err)
	}

	// create role
	p.Log.Infof("Create iam role %s...", devPodRoleName)
	roleOutput, err := iamClient.CreateRole(ctx, &iam.CreateRoleInput{
		RoleName: &devPodRoleName,
		AssumeRolePolicyDocument: options.Ptr(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`),
	})
	if err != nil {
		_, _ = iamClient.DeletePolicy(ctx, &iam.DeletePolicyInput{PolicyArn: policyOutput.Policy.Arn})
		return "", fmt.Errorf("create iam role: %w", err)
	}

	// attach policy
	p.Log.Infof("Attach iam policy %s to role %s...", devPodPolicyName, devPodRoleName)
	_, err = iamClient.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
		PolicyArn: policyOutput.Policy.Arn,
		RoleName:  &devPodRoleName,
	})
	if err != nil {
		_, _ = iamClient.DeletePolicy(ctx, &iam.DeletePolicyInput{PolicyArn: policyOutput.Policy.Arn})
		_, _ = iamClient.DeleteRole(ctx, &iam.DeleteRoleInput{RoleName: &devPodRoleName})
		return "", fmt.Errorf("attach iam policy to role: %w", err)
	}

	return *roleOutput.Role.Arn, nil
}
