package aws

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	"github.com/awesome-foundation/cfnpeek/internal/model"
)

// CloudFormationAPI abstracts the CloudFormation API calls for testability.
type CloudFormationAPI interface {
	DescribeStacks(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error)
	DescribeStackEvents(ctx context.Context, params *cloudformation.DescribeStackEventsInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStackEventsOutput, error)
	ListStacks(ctx context.Context, params *cloudformation.ListStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ListStacksOutput, error)
	ListStackResources(ctx context.Context, params *cloudformation.ListStackResourcesInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ListStackResourcesOutput, error)
	ListExports(ctx context.Context, params *cloudformation.ListExportsInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ListExportsOutput, error)
}

// Client wraps an AWS CloudFormation API client.
type Client struct {
	api CloudFormationAPI
}

// NewClient creates a Client configured with the given region and profile.
func NewClient(ctx context.Context, region, profile string) (*Client, error) {
	var opts []func(*config.LoadOptions) error
	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}
	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}

	return &Client{api: cloudformation.NewFromConfig(cfg)}, nil
}

// NewClientFromAPI creates a Client from an existing API implementation (for testing).
func NewClientFromAPI(api CloudFormationAPI) *Client {
	return &Client{api: api}
}

// activeStatuses are the stack statuses we show by default (excludes DELETE_COMPLETE).
var activeStatuses = []cftypes.StackStatus{
	cftypes.StackStatusCreateInProgress,
	cftypes.StackStatusCreateFailed,
	cftypes.StackStatusCreateComplete,
	cftypes.StackStatusRollbackInProgress,
	cftypes.StackStatusRollbackFailed,
	cftypes.StackStatusRollbackComplete,
	cftypes.StackStatusDeleteInProgress,
	cftypes.StackStatusDeleteFailed,
	cftypes.StackStatusUpdateInProgress,
	cftypes.StackStatusUpdateCompleteCleanupInProgress,
	cftypes.StackStatusUpdateComplete,
	cftypes.StackStatusUpdateFailed,
	cftypes.StackStatusUpdateRollbackInProgress,
	cftypes.StackStatusUpdateRollbackFailed,
	cftypes.StackStatusUpdateRollbackCompleteCleanupInProgress,
	cftypes.StackStatusUpdateRollbackComplete,
	cftypes.StackStatusReviewInProgress,
	cftypes.StackStatusImportInProgress,
	cftypes.StackStatusImportComplete,
	cftypes.StackStatusImportRollbackInProgress,
	cftypes.StackStatusImportRollbackFailed,
	cftypes.StackStatusImportRollbackComplete,
}

// FetchStacks returns a list of all active stacks in the region.
func (c *Client) FetchStacks(ctx context.Context) (*model.StackList, error) {
	var stacks []model.StackSummary
	var nextToken *string

	for {
		out, err := c.api.ListStacks(ctx, &cloudformation.ListStacksInput{
			StackStatusFilter: activeStatuses,
			NextToken:         nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("listing stacks: %w", err)
		}

		for _, s := range out.StackSummaries {
			createdAt := ""
			if s.CreationTime != nil {
				createdAt = s.CreationTime.Format("2006-01-02T15:04:05Z")
			}
			updatedAt := ""
			if s.LastUpdatedTime != nil {
				updatedAt = s.LastUpdatedTime.Format("2006-01-02T15:04:05Z")
			}
			stacks = append(stacks, model.StackSummary{
				StackName:   aws.ToString(s.StackName),
				StackID:     aws.ToString(s.StackId),
				Status:      string(s.StackStatus),
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
				Description: aws.ToString(s.TemplateDescription),
			})
		}

		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}

	sort.Slice(stacks, func(i, j int) bool {
		return stacks[i].StackName < stacks[j].StackName
	})

	return &model.StackList{Stacks: stacks}, nil
}

// FetchResources returns all resources for a stack, handling pagination.
func (c *Client) FetchResources(ctx context.Context, stackName string) ([]model.Resource, error) {
	var resources []model.Resource
	var nextToken *string

	for {
		out, err := c.api.ListStackResources(ctx, &cloudformation.ListStackResourcesInput{
			StackName: aws.String(stackName),
			NextToken: nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("listing stack resources: %w", err)
		}

		for _, r := range out.StackResourceSummaries {
			lastUpdated := ""
			if r.LastUpdatedTimestamp != nil {
				lastUpdated = r.LastUpdatedTimestamp.Format("2006-01-02T15:04:05Z")
			}
			resources = append(resources, model.Resource{
				LogicalID:   aws.ToString(r.LogicalResourceId),
				PhysicalID:  aws.ToString(r.PhysicalResourceId),
				Type:        aws.ToString(r.ResourceType),
				Status:      string(r.ResourceStatus),
				LastUpdated: lastUpdated,
			})
		}

		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}

	return resources, nil
}

// FetchStack returns stack metadata and outputs.
func (c *Client) FetchStack(ctx context.Context, stackName string) (string, string, string, []model.Output, error) {
	out, err := c.api.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return "", "", "", nil, fmt.Errorf("describing stack: %w", err)
	}

	if len(out.Stacks) == 0 {
		return "", "", "", nil, fmt.Errorf("stack %q not found", stackName)
	}

	stack := out.Stacks[0]
	var outputs []model.Output
	for _, o := range stack.Outputs {
		outputs = append(outputs, model.Output{
			Key:         aws.ToString(o.OutputKey),
			Value:       aws.ToString(o.OutputValue),
			Description: aws.ToString(o.Description),
			ExportName:  aws.ToString(o.ExportName),
		})
	}

	return aws.ToString(stack.StackName), aws.ToString(stack.StackId), string(stack.StackStatus), outputs, nil
}

// FetchExports returns all exports belonging to the given stack.
// ListExports is account-wide, so we filter by the stack's ID.
func (c *Client) FetchExports(ctx context.Context, stackID string) ([]model.Export, error) {
	var exports []model.Export
	var nextToken *string

	for {
		out, err := c.api.ListExports(ctx, &cloudformation.ListExportsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("listing exports: %w", err)
		}

		for _, e := range out.Exports {
			if aws.ToString(e.ExportingStackId) == stackID {
				exports = append(exports, model.Export{
					Name:  aws.ToString(e.Name),
					Value: aws.ToString(e.Value),
				})
			}
		}

		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}

	return exports, nil
}

// FetchEvents returns stack events sorted ascending by timestamp (oldest first).
// AWS returns events newest-first, so we reverse after collecting all pages.
// If limit > 0, only the last N events (by time) are returned.
func (c *Client) FetchEvents(ctx context.Context, stackName string, limit int) (*model.StackEvents, error) {
	var events []model.StackEvent
	var nextToken *string

	for {
		out, err := c.api.DescribeStackEvents(ctx, &cloudformation.DescribeStackEventsInput{
			StackName: aws.String(stackName),
			NextToken: nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("describing stack events: %w", err)
		}

		for _, e := range out.StackEvents {
			ts := ""
			if e.Timestamp != nil {
				ts = e.Timestamp.Format("2006-01-02T15:04:05Z")
			}
			events = append(events, model.StackEvent{
				Timestamp:    ts,
				LogicalID:    aws.ToString(e.LogicalResourceId),
				Status:       string(e.ResourceStatus),
				StatusReason: aws.ToString(e.ResourceStatusReason),
				ResourceType: aws.ToString(e.ResourceType),
				PhysicalID:   aws.ToString(e.PhysicalResourceId),
			})
		}

		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}

	// AWS returns newest-first; reverse to get oldest-first.
	for i, j := 0, len(events)-1; i < j; i, j = i+1, j-1 {
		events[i], events[j] = events[j], events[i]
	}

	// Apply limit: keep the last N events (already sorted ascending, so tail).
	if limit > 0 && len(events) > limit {
		events = events[len(events)-limit:]
	}

	return &model.StackEvents{
		StackName: stackName,
		Events:    events,
	}, nil
}

// FetchStackInfo fetches all requested sections for a stack.
func (c *Client) FetchStackInfo(ctx context.Context, stackName string, wantResources, wantOutputs, wantExports bool) (*model.StackInfo, error) {
	name, id, status, outputs, err := c.FetchStack(ctx, stackName)
	if err != nil {
		return nil, err
	}

	info := &model.StackInfo{
		StackName: name,
		StackID:   id,
		Status:    status,
	}

	if wantOutputs {
		sort.Slice(outputs, func(i, j int) bool {
			return outputs[i].Key < outputs[j].Key
		})
		info.Outputs = outputs
	}

	if wantResources {
		resources, err := c.FetchResources(ctx, stackName)
		if err != nil {
			return nil, err
		}
		sort.Slice(resources, func(i, j int) bool {
			return resources[i].LogicalID < resources[j].LogicalID
		})
		info.Resources = resources
	}

	if wantExports {
		exports, err := c.FetchExports(ctx, id)
		if err != nil {
			return nil, err
		}
		sort.Slice(exports, func(i, j int) bool {
			return exports[i].Name < exports[j].Name
		})
		info.Exports = exports
	}

	return info, nil
}

// FormatError returns a user-friendly error message for common AWS errors.
func FormatError(err error) string {
	msg := err.Error()

	if strings.Contains(msg, "does not exist") {
		return "Stack not found. Check the stack name/ARN and region."
	}
	if strings.Contains(msg, "ExpiredToken") || strings.Contains(msg, "ExpiredTokenException") {
		return "AWS credentials have expired. Refresh your credentials and try again."
	}
	if strings.Contains(msg, "NoCredentialProviders") || strings.Contains(msg, "no EC2 IMDS") {
		return "No AWS credentials found. Set AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY, use --profile, or configure AWS SSO."
	}
	if strings.Contains(msg, "AccessDenied") || strings.Contains(msg, "is not authorized") {
		return "Access denied. Check your IAM permissions for CloudFormation read access."
	}
	if strings.Contains(msg, "could not find region") || strings.Contains(msg, "MissingRegion") {
		return "No AWS region configured. Use --region, set AWS_REGION, or configure a default region."
	}

	// Check for ValidationError (common for bad stack names)
	var ve *cftypes.StackNotFoundException
	_ = ve // type assertion happens via errors.As in callers if needed

	return fmt.Sprintf("AWS error: %s", msg)
}
