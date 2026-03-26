package aws_test

import (
	"context"
	"testing"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	cfnaws "github.com/awesome-foundation/cfnpeek/internal/aws"
)

type mockAPI struct {
	describeStacksOutput      *cloudformation.DescribeStacksOutput
	describeStackEventsPages  []*cloudformation.DescribeStackEventsOutput
	listStacksPages           []*cloudformation.ListStacksOutput
	listStackResourcesPages   []*cloudformation.ListStackResourcesOutput
	listExportsPages          []*cloudformation.ListExportsOutput
	describeErr               error
	describeEventsErr         error
	listStacksErr             error
	listResourcesErr          error
	listExportsErr            error
	stacksPage                int
	resourcePage              int
	exportPage                int
	eventsPage                int
}

func (m *mockAPI) ListStacks(_ context.Context, _ *cloudformation.ListStacksInput, _ ...func(*cloudformation.Options)) (*cloudformation.ListStacksOutput, error) {
	if m.listStacksErr != nil {
		return nil, m.listStacksErr
	}
	if len(m.listStacksPages) == 0 {
		return &cloudformation.ListStacksOutput{}, nil
	}
	page := m.listStacksPages[m.stacksPage]
	m.stacksPage++
	return page, nil
}

func (m *mockAPI) DescribeStacks(_ context.Context, _ *cloudformation.DescribeStacksInput, _ ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
	if m.describeErr != nil {
		return nil, m.describeErr
	}
	return m.describeStacksOutput, nil
}

func (m *mockAPI) DescribeStackEvents(_ context.Context, _ *cloudformation.DescribeStackEventsInput, _ ...func(*cloudformation.Options)) (*cloudformation.DescribeStackEventsOutput, error) {
	if m.describeEventsErr != nil {
		return nil, m.describeEventsErr
	}
	if len(m.describeStackEventsPages) == 0 {
		return &cloudformation.DescribeStackEventsOutput{}, nil
	}
	page := m.describeStackEventsPages[m.eventsPage]
	m.eventsPage++
	return page, nil
}

func (m *mockAPI) ListStackResources(_ context.Context, _ *cloudformation.ListStackResourcesInput, _ ...func(*cloudformation.Options)) (*cloudformation.ListStackResourcesOutput, error) {
	if m.listResourcesErr != nil {
		return nil, m.listResourcesErr
	}
	page := m.listStackResourcesPages[m.resourcePage]
	m.resourcePage++
	return page, nil
}

func (m *mockAPI) ListExports(_ context.Context, _ *cloudformation.ListExportsInput, _ ...func(*cloudformation.Options)) (*cloudformation.ListExportsOutput, error) {
	if m.listExportsErr != nil {
		return nil, m.listExportsErr
	}
	page := m.listExportsPages[m.exportPage]
	m.exportPage++
	return page, nil
}

func TestFetchStackInfo(t *testing.T) {
	ts := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	stackID := "arn:aws:cloudformation:us-east-1:123456789:stack/test-stack/guid"

	mock := &mockAPI{
		describeStacksOutput: &cloudformation.DescribeStacksOutput{
			Stacks: []cftypes.Stack{
				{
					StackName:   awssdk.String("test-stack"),
					StackId:     awssdk.String(stackID),
					StackStatus: cftypes.StackStatusCreateComplete,
					Outputs: []cftypes.Output{
						{
							OutputKey:   awssdk.String("BucketArn"),
							OutputValue: awssdk.String("arn:aws:s3:::my-bucket"),
							ExportName:  awssdk.String("test-stack-BucketArn"),
						},
					},
				},
			},
		},
		listStackResourcesPages: []*cloudformation.ListStackResourcesOutput{
			{
				StackResourceSummaries: []cftypes.StackResourceSummary{
					{
						LogicalResourceId:    awssdk.String("MyBucket"),
						PhysicalResourceId:   awssdk.String("my-bucket-abc123"),
						ResourceType:         awssdk.String("AWS::S3::Bucket"),
						ResourceStatus:       cftypes.ResourceStatusCreateComplete,
						LastUpdatedTimestamp: &ts,
					},
				},
			},
		},
		listExportsPages: []*cloudformation.ListExportsOutput{
			{
				Exports: []cftypes.Export{
					{
						Name:             awssdk.String("test-stack-BucketArn"),
						Value:            awssdk.String("arn:aws:s3:::my-bucket"),
						ExportingStackId: awssdk.String(stackID),
					},
					{
						Name:             awssdk.String("other-stack-Export"),
						Value:            awssdk.String("other-value"),
						ExportingStackId: awssdk.String("arn:aws:cloudformation:us-east-1:123456789:stack/other-stack/guid"),
					},
				},
			},
		},
	}

	client := cfnaws.NewClientFromAPI(mock)
	info, err := client.FetchStackInfo(context.Background(), "test-stack", true, true, true)
	if err != nil {
		t.Fatal(err)
	}

	if info.StackName != "test-stack" {
		t.Errorf("expected stack name 'test-stack', got %q", info.StackName)
	}
	if info.Status != "CREATE_COMPLETE" {
		t.Errorf("expected status CREATE_COMPLETE, got %q", info.Status)
	}
	if len(info.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(info.Resources))
	}
	if info.Resources[0].Type != "AWS::S3::Bucket" {
		t.Errorf("expected resource type AWS::S3::Bucket, got %q", info.Resources[0].Type)
	}
	if len(info.Outputs) != 1 {
		t.Fatalf("expected 1 output, got %d", len(info.Outputs))
	}
	// Exports should be filtered to only this stack
	if len(info.Exports) != 1 {
		t.Fatalf("expected 1 export (filtered), got %d", len(info.Exports))
	}
	if info.Exports[0].Name != "test-stack-BucketArn" {
		t.Errorf("expected export name 'test-stack-BucketArn', got %q", info.Exports[0].Name)
	}
}

func TestFetchEvents(t *testing.T) {
	t1 := time.Date(2026, 1, 15, 10, 30, 5, 0, time.UTC)
	t2 := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC) // older

	mock := &mockAPI{
		describeStackEventsPages: []*cloudformation.DescribeStackEventsOutput{
			{
				// AWS returns newest first
				StackEvents: []cftypes.StackEvent{
					{
						Timestamp:            &t1,
						LogicalResourceId:    awssdk.String("Subnet"),
						ResourceStatus:       cftypes.ResourceStatusCreateFailed,
						ResourceStatusReason: awssdk.String("Resource limit exceeded"),
						ResourceType:         awssdk.String("AWS::EC2::Subnet"),
						PhysicalResourceId:   awssdk.String(""),
					},
					{
						Timestamp:          &t2,
						LogicalResourceId:  awssdk.String("VPC"),
						ResourceStatus:     cftypes.ResourceStatusCreateComplete,
						ResourceType:       awssdk.String("AWS::EC2::VPC"),
						PhysicalResourceId: awssdk.String("vpc-abc123"),
					},
				},
			},
		},
	}

	client := cfnaws.NewClientFromAPI(mock)

	t.Run("all events sorted ascending", func(t *testing.T) {
		result, err := client.FetchEvents(context.Background(), "test-stack", 0)
		if err != nil {
			t.Fatal(err)
		}
		if result.StackName != "test-stack" {
			t.Errorf("expected stack name 'test-stack', got %q", result.StackName)
		}
		if len(result.Events) != 2 {
			t.Fatalf("expected 2 events, got %d", len(result.Events))
		}
		// After reversal the older event (VPC) should be first
		if result.Events[0].LogicalID != "VPC" {
			t.Errorf("expected first event LogicalID 'VPC', got %q", result.Events[0].LogicalID)
		}
		if result.Events[1].LogicalID != "Subnet" {
			t.Errorf("expected second event LogicalID 'Subnet', got %q", result.Events[1].LogicalID)
		}
		if result.Events[1].StatusReason != "Resource limit exceeded" {
			t.Errorf("expected status reason 'Resource limit exceeded', got %q", result.Events[1].StatusReason)
		}
	})
}

func TestFetchEventsLimit(t *testing.T) {
	t1 := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
	t2 := time.Date(2026, 1, 15, 10, 30, 5, 0, time.UTC)
	t3 := time.Date(2026, 1, 15, 10, 30, 10, 0, time.UTC)

	mock := &mockAPI{
		describeStackEventsPages: []*cloudformation.DescribeStackEventsOutput{
			{
				// AWS returns newest first
				StackEvents: []cftypes.StackEvent{
					{Timestamp: &t3, LogicalResourceId: awssdk.String("C"), ResourceType: awssdk.String("AWS::CloudFormation::Stack"), ResourceStatus: cftypes.ResourceStatusCreateComplete},
					{Timestamp: &t2, LogicalResourceId: awssdk.String("B"), ResourceType: awssdk.String("AWS::CloudFormation::Stack"), ResourceStatus: cftypes.ResourceStatusCreateComplete},
					{Timestamp: &t1, LogicalResourceId: awssdk.String("A"), ResourceType: awssdk.String("AWS::CloudFormation::Stack"), ResourceStatus: cftypes.ResourceStatusCreateComplete},
				},
			},
		},
	}

	client := cfnaws.NewClientFromAPI(mock)

	result, err := client.FetchEvents(context.Background(), "test-stack", 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Events) != 2 {
		t.Fatalf("expected 2 events after limit, got %d", len(result.Events))
	}
	// Should keep the last 2 (B and C, i.e. the most recent)
	if result.Events[0].LogicalID != "B" {
		t.Errorf("expected first event 'B', got %q", result.Events[0].LogicalID)
	}
	if result.Events[1].LogicalID != "C" {
		t.Errorf("expected second event 'C', got %q", result.Events[1].LogicalID)
	}
}

func TestFetchStackInfoSectionsFiltered(t *testing.T) {
	stackID := "arn:aws:cloudformation:us-east-1:123456789:stack/test-stack/guid"

	mock := &mockAPI{
		describeStacksOutput: &cloudformation.DescribeStacksOutput{
			Stacks: []cftypes.Stack{
				{
					StackName:   awssdk.String("test-stack"),
					StackId:     awssdk.String(stackID),
					StackStatus: cftypes.StackStatusCreateComplete,
				},
			},
		},
	}

	client := cfnaws.NewClientFromAPI(mock)

	// Only outputs requested, no resources or exports calls
	info, err := client.FetchStackInfo(context.Background(), "test-stack", false, true, false)
	if err != nil {
		t.Fatal(err)
	}
	if info.Resources != nil {
		t.Error("expected nil resources when not requested")
	}
	if info.Exports != nil {
		t.Error("expected nil exports when not requested")
	}
}
