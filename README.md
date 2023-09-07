# ECS Provider for DevPod

## Getting started

The provider is available for auto-installation using 

```sh
devpod provider add ecs
devpod provider use ecs
```

Follow the on-screen instructions to complete the setup.

Needed variables will be:

- **CLUSTER_ID**: ECS Cluster ID either as ARN or ID
- **SUBNET_ID**: ECS Subnet ID either as ARN or ID to run the tasks in. This can either be a private subnet with a NAT Gateway or a Public Subnet. Depending on the type of the subnet you will need to set ASSIGN_PUBLIC_IP accordingly
- **ASSIGN_PUBLIC_IP**: If the task should get a public ip assigned. For public subnets specify ENABLED and for private subnets specify DISABLED.
- **TASK_ROLE_ARN**: ECS Task Role ARN to use for the task definition with IAM permissions required for ECS Exec. For more information take a look at https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-iam-roles.html
- **EXECUTION_ROLE_ARN**: ECS Execution Role ARN to use for the task definition. For more information take a look at https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_execution_IAM_role.html
- **LAUNCH_TYPE**: ECS Task Launch Type, which can be either FARGATE, ECS or EXTERNAL. If the LAUNCH_TYPE is FARGATE, stopping workspaces is not supported.

The provider will inherit the login information from `aws cli` or you can
specify in your environment, or in the provider options, the `AWS_ACCESS_KEY_ID=`
and `AWS_SECRET_ACCESS_KEY=`

### Creating your first DevPod env with aws

After the initial setup, just use:

```sh
devpod up .
```

You'll need to wait for the task and environment setup.
