# ECS Provider for DevPod

[![Join us on Slack!](docs/static/media/slack.svg)](https://slack.loft.sh/) [![Open in DevPod!](https://devpod.sh/assets/open-in-devpod.svg)](https://devpod.sh/open#https://github.com/loft-sh/devpod-provider-ecs)

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

The provider will inherit the login information from `aws cli` or you can
specify in your environment the `AWS_ACCESS_KEY_ID=`
and `AWS_SECRET_ACCESS_KEY=` variables.

### Creating your first DevPod env with aws

After the initial setup, just use:

```sh
devpod up .
```

You'll need to wait for the task and environment setup.
