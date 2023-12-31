// Code generated by smithy-go-codegen DO NOT EDIT.

package ssm

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	internalauth "github.com/aws/aws-sdk-go-v2/internal/auth"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Retrieves the details of a maintenance window task. For maintenance window
// tasks without a specified target, you can't supply values for --max-errors and
// --max-concurrency . Instead, the system inserts a placeholder value of 1 , which
// may be reported in the response to this command. These values don't affect the
// running of your task and can be ignored. To retrieve a list of tasks in a
// maintenance window, instead use the DescribeMaintenanceWindowTasks command.
func (c *Client) GetMaintenanceWindowTask(ctx context.Context, params *GetMaintenanceWindowTaskInput, optFns ...func(*Options)) (*GetMaintenanceWindowTaskOutput, error) {
	if params == nil {
		params = &GetMaintenanceWindowTaskInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "GetMaintenanceWindowTask", params, optFns, c.addOperationGetMaintenanceWindowTaskMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*GetMaintenanceWindowTaskOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type GetMaintenanceWindowTaskInput struct {

	// The maintenance window ID that includes the task to retrieve.
	//
	// This member is required.
	WindowId *string

	// The maintenance window task ID to retrieve.
	//
	// This member is required.
	WindowTaskId *string

	noSmithyDocumentSerde
}

type GetMaintenanceWindowTaskOutput struct {

	// The details for the CloudWatch alarm you applied to your maintenance window
	// task.
	AlarmConfiguration *types.AlarmConfiguration

	// The action to take on tasks when the maintenance window cutoff time is reached.
	// CONTINUE_TASK means that tasks continue to run. For Automation, Lambda, Step
	// Functions tasks, CANCEL_TASK means that currently running task invocations
	// continue, but no new task invocations are started. For Run Command tasks,
	// CANCEL_TASK means the system attempts to stop the task by sending a
	// CancelCommand operation.
	CutoffBehavior types.MaintenanceWindowTaskCutoffBehavior

	// The retrieved task description.
	Description *string

	// The location in Amazon Simple Storage Service (Amazon S3) where the task
	// results are logged. LoggingInfo has been deprecated. To specify an Amazon
	// Simple Storage Service (Amazon S3) bucket to contain logs, instead use the
	// OutputS3BucketName and OutputS3KeyPrefix options in the TaskInvocationParameters
	// structure. For information about how Amazon Web Services Systems Manager handles
	// these options for the supported maintenance window task types, see
	// MaintenanceWindowTaskInvocationParameters .
	LoggingInfo *types.LoggingInfo

	// The maximum number of targets allowed to run this task in parallel. For
	// maintenance window tasks without a target specified, you can't supply a value
	// for this option. Instead, the system inserts a placeholder value of 1 , which
	// may be reported in the response to this command. This value doesn't affect the
	// running of your task and can be ignored.
	MaxConcurrency *string

	// The maximum number of errors allowed before the task stops being scheduled. For
	// maintenance window tasks without a target specified, you can't supply a value
	// for this option. Instead, the system inserts a placeholder value of 1 , which
	// may be reported in the response to this command. This value doesn't affect the
	// running of your task and can be ignored.
	MaxErrors *string

	// The retrieved task name.
	Name *string

	// The priority of the task when it runs. The lower the number, the higher the
	// priority. Tasks that have the same priority are scheduled in parallel.
	Priority int32

	// The Amazon Resource Name (ARN) of the Identity and Access Management (IAM)
	// service role to use to publish Amazon Simple Notification Service (Amazon SNS)
	// notifications for maintenance window Run Command tasks.
	ServiceRoleArn *string

	// The targets where the task should run.
	Targets []types.Target

	// The resource that the task used during execution. For RUN_COMMAND and AUTOMATION
	// task types, the value of TaskArn is the SSM document name/ARN. For LAMBDA
	// tasks, the value is the function name/ARN. For STEP_FUNCTIONS tasks, the value
	// is the state machine ARN.
	TaskArn *string

	// The parameters to pass to the task when it runs.
	TaskInvocationParameters *types.MaintenanceWindowTaskInvocationParameters

	// The parameters to pass to the task when it runs. TaskParameters has been
	// deprecated. To specify parameters to pass to a task when it runs, instead use
	// the Parameters option in the TaskInvocationParameters structure. For
	// information about how Systems Manager handles these options for the supported
	// maintenance window task types, see MaintenanceWindowTaskInvocationParameters .
	TaskParameters map[string]types.MaintenanceWindowTaskParameterValueExpression

	// The type of task to run.
	TaskType types.MaintenanceWindowTaskType

	// The retrieved maintenance window ID.
	WindowId *string

	// The retrieved maintenance window task ID.
	WindowTaskId *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationGetMaintenanceWindowTaskMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpGetMaintenanceWindowTask{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpGetMaintenanceWindowTask{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addGetMaintenanceWindowTaskResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addOpGetMaintenanceWindowTaskValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opGetMaintenanceWindowTask(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addendpointDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opGetMaintenanceWindowTask(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "ssm",
		OperationName: "GetMaintenanceWindowTask",
	}
}

type opGetMaintenanceWindowTaskResolveEndpointMiddleware struct {
	EndpointResolver EndpointResolverV2
	BuiltInResolver  builtInParameterResolver
}

func (*opGetMaintenanceWindowTaskResolveEndpointMiddleware) ID() string {
	return "ResolveEndpointV2"
}

func (m *opGetMaintenanceWindowTaskResolveEndpointMiddleware) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	if awsmiddleware.GetRequiresLegacyEndpoints(ctx) {
		return next.HandleSerialize(ctx, in)
	}

	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown transport type %T", in.Request)
	}

	if m.EndpointResolver == nil {
		return out, metadata, fmt.Errorf("expected endpoint resolver to not be nil")
	}

	params := EndpointParameters{}

	m.BuiltInResolver.ResolveBuiltIns(&params)

	var resolvedEndpoint smithyendpoints.Endpoint
	resolvedEndpoint, err = m.EndpointResolver.ResolveEndpoint(ctx, params)
	if err != nil {
		return out, metadata, fmt.Errorf("failed to resolve service endpoint, %w", err)
	}

	req.URL = &resolvedEndpoint.URI

	for k := range resolvedEndpoint.Headers {
		req.Header.Set(
			k,
			resolvedEndpoint.Headers.Get(k),
		)
	}

	authSchemes, err := internalauth.GetAuthenticationSchemes(&resolvedEndpoint.Properties)
	if err != nil {
		var nfe *internalauth.NoAuthenticationSchemesFoundError
		if errors.As(err, &nfe) {
			// if no auth scheme is found, default to sigv4
			signingName := "ssm"
			signingRegion := m.BuiltInResolver.(*builtInResolver).Region
			ctx = awsmiddleware.SetSigningName(ctx, signingName)
			ctx = awsmiddleware.SetSigningRegion(ctx, signingRegion)

		}
		var ue *internalauth.UnSupportedAuthenticationSchemeSpecifiedError
		if errors.As(err, &ue) {
			return out, metadata, fmt.Errorf(
				"This operation requests signer version(s) %v but the client only supports %v",
				ue.UnsupportedSchemes,
				internalauth.SupportedSchemes,
			)
		}
	}

	for _, authScheme := range authSchemes {
		switch authScheme.(type) {
		case *internalauth.AuthenticationSchemeV4:
			v4Scheme, _ := authScheme.(*internalauth.AuthenticationSchemeV4)
			var signingName, signingRegion string
			if v4Scheme.SigningName == nil {
				signingName = "ssm"
			} else {
				signingName = *v4Scheme.SigningName
			}
			if v4Scheme.SigningRegion == nil {
				signingRegion = m.BuiltInResolver.(*builtInResolver).Region
			} else {
				signingRegion = *v4Scheme.SigningRegion
			}
			if v4Scheme.DisableDoubleEncoding != nil {
				// The signer sets an equivalent value at client initialization time.
				// Setting this context value will cause the signer to extract it
				// and override the value set at client initialization time.
				ctx = internalauth.SetDisableDoubleEncoding(ctx, *v4Scheme.DisableDoubleEncoding)
			}
			ctx = awsmiddleware.SetSigningName(ctx, signingName)
			ctx = awsmiddleware.SetSigningRegion(ctx, signingRegion)
			break
		case *internalauth.AuthenticationSchemeV4A:
			v4aScheme, _ := authScheme.(*internalauth.AuthenticationSchemeV4A)
			if v4aScheme.SigningName == nil {
				v4aScheme.SigningName = aws.String("ssm")
			}
			if v4aScheme.DisableDoubleEncoding != nil {
				// The signer sets an equivalent value at client initialization time.
				// Setting this context value will cause the signer to extract it
				// and override the value set at client initialization time.
				ctx = internalauth.SetDisableDoubleEncoding(ctx, *v4aScheme.DisableDoubleEncoding)
			}
			ctx = awsmiddleware.SetSigningName(ctx, *v4aScheme.SigningName)
			ctx = awsmiddleware.SetSigningRegion(ctx, v4aScheme.SigningRegionSet[0])
			break
		case *internalauth.AuthenticationSchemeNone:
			break
		}
	}

	return next.HandleSerialize(ctx, in)
}

func addGetMaintenanceWindowTaskResolveEndpointMiddleware(stack *middleware.Stack, options Options) error {
	return stack.Serialize.Insert(&opGetMaintenanceWindowTaskResolveEndpointMiddleware{
		EndpointResolver: options.EndpointResolverV2,
		BuiltInResolver: &builtInResolver{
			Region:       options.Region,
			UseDualStack: options.EndpointOptions.UseDualStackEndpoint,
			UseFIPS:      options.EndpointOptions.UseFIPSEndpoint,
			Endpoint:     options.BaseEndpoint,
		},
	}, "ResolveEndpoint", middleware.After)
}
