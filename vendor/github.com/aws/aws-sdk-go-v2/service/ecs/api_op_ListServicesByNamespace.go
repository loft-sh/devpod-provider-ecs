// Code generated by smithy-go-codegen DO NOT EDIT.

package ecs

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	internalauth "github.com/aws/aws-sdk-go-v2/internal/auth"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// This operation lists all of the services that are associated with a Cloud Map
// namespace. This list might include services in different clusters. In contrast,
// ListServices can only list services in one cluster at a time. If you need to
// filter the list of services in a single cluster by various parameters, use
// ListServices . For more information, see Service Connect (https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service-connect.html)
// in the Amazon Elastic Container Service Developer Guide.
func (c *Client) ListServicesByNamespace(ctx context.Context, params *ListServicesByNamespaceInput, optFns ...func(*Options)) (*ListServicesByNamespaceOutput, error) {
	if params == nil {
		params = &ListServicesByNamespaceInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ListServicesByNamespace", params, optFns, c.addOperationListServicesByNamespaceMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ListServicesByNamespaceOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ListServicesByNamespaceInput struct {

	// The namespace name or full Amazon Resource Name (ARN) of the Cloud Map
	// namespace to list the services in. Tasks that run in a namespace can use short
	// names to connect to services in the namespace. Tasks can connect to services
	// across all of the clusters in the namespace. Tasks connect through a managed
	// proxy container that collects logs and metrics for increased visibility. Only
	// the tasks that Amazon ECS services create are supported with Service Connect.
	// For more information, see Service Connect (https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service-connect.html)
	// in the Amazon Elastic Container Service Developer Guide.
	//
	// This member is required.
	Namespace *string

	// The maximum number of service results that ListServicesByNamespace returns in
	// paginated output. When this parameter is used, ListServicesByNamespace only
	// returns maxResults results in a single page along with a nextToken response
	// element. The remaining results of the initial request can be seen by sending
	// another ListServicesByNamespace request with the returned nextToken value. This
	// value can be between 1 and 100. If this parameter isn't used, then
	// ListServicesByNamespace returns up to 10 results and a nextToken value if
	// applicable.
	MaxResults *int32

	// The nextToken value that's returned from a ListServicesByNamespace request. It
	// indicates that more results are available to fulfill the request and further
	// calls are needed. If maxResults is returned, it is possible the number of
	// results is less than maxResults .
	NextToken *string

	noSmithyDocumentSerde
}

type ListServicesByNamespaceOutput struct {

	// The nextToken value to include in a future ListServicesByNamespace request.
	// When the results of a ListServicesByNamespace request exceed maxResults , this
	// value can be used to retrieve the next page of results. When there are no more
	// results to return, this value is null .
	NextToken *string

	// The list of full ARN entries for each service that's associated with the
	// specified namespace.
	ServiceArns []string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationListServicesByNamespaceMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpListServicesByNamespace{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpListServicesByNamespace{}, middleware.After)
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
	if err = addListServicesByNamespaceResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addOpListServicesByNamespaceValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opListServicesByNamespace(options.Region), middleware.Before); err != nil {
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

// ListServicesByNamespaceAPIClient is a client that implements the
// ListServicesByNamespace operation.
type ListServicesByNamespaceAPIClient interface {
	ListServicesByNamespace(context.Context, *ListServicesByNamespaceInput, ...func(*Options)) (*ListServicesByNamespaceOutput, error)
}

var _ ListServicesByNamespaceAPIClient = (*Client)(nil)

// ListServicesByNamespacePaginatorOptions is the paginator options for
// ListServicesByNamespace
type ListServicesByNamespacePaginatorOptions struct {
	// The maximum number of service results that ListServicesByNamespace returns in
	// paginated output. When this parameter is used, ListServicesByNamespace only
	// returns maxResults results in a single page along with a nextToken response
	// element. The remaining results of the initial request can be seen by sending
	// another ListServicesByNamespace request with the returned nextToken value. This
	// value can be between 1 and 100. If this parameter isn't used, then
	// ListServicesByNamespace returns up to 10 results and a nextToken value if
	// applicable.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// ListServicesByNamespacePaginator is a paginator for ListServicesByNamespace
type ListServicesByNamespacePaginator struct {
	options   ListServicesByNamespacePaginatorOptions
	client    ListServicesByNamespaceAPIClient
	params    *ListServicesByNamespaceInput
	nextToken *string
	firstPage bool
}

// NewListServicesByNamespacePaginator returns a new
// ListServicesByNamespacePaginator
func NewListServicesByNamespacePaginator(client ListServicesByNamespaceAPIClient, params *ListServicesByNamespaceInput, optFns ...func(*ListServicesByNamespacePaginatorOptions)) *ListServicesByNamespacePaginator {
	if params == nil {
		params = &ListServicesByNamespaceInput{}
	}

	options := ListServicesByNamespacePaginatorOptions{}
	if params.MaxResults != nil {
		options.Limit = *params.MaxResults
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListServicesByNamespacePaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.NextToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *ListServicesByNamespacePaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next ListServicesByNamespace page.
func (p *ListServicesByNamespacePaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*ListServicesByNamespaceOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxResults = limit

	result, err := p.client.ListServicesByNamespace(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextToken

	if p.options.StopOnDuplicateToken &&
		prevToken != nil &&
		p.nextToken != nil &&
		*prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

func newServiceMetadataMiddleware_opListServicesByNamespace(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "ecs",
		OperationName: "ListServicesByNamespace",
	}
}

type opListServicesByNamespaceResolveEndpointMiddleware struct {
	EndpointResolver EndpointResolverV2
	BuiltInResolver  builtInParameterResolver
}

func (*opListServicesByNamespaceResolveEndpointMiddleware) ID() string {
	return "ResolveEndpointV2"
}

func (m *opListServicesByNamespaceResolveEndpointMiddleware) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
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
			signingName := "ecs"
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
				signingName = "ecs"
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
				v4aScheme.SigningName = aws.String("ecs")
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

func addListServicesByNamespaceResolveEndpointMiddleware(stack *middleware.Stack, options Options) error {
	return stack.Serialize.Insert(&opListServicesByNamespaceResolveEndpointMiddleware{
		EndpointResolver: options.EndpointResolverV2,
		BuiltInResolver: &builtInResolver{
			Region:       options.Region,
			UseDualStack: options.EndpointOptions.UseDualStackEndpoint,
			UseFIPS:      options.EndpointOptions.UseFIPSEndpoint,
			Endpoint:     options.BaseEndpoint,
		},
	}, "ResolveEndpoint", middleware.After)
}
