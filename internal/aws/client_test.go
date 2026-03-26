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
	describeStacksOutput    *cloudformation.DescribeStacksOutput
	listStacksPages         []*cloudformation.ListStacksOutput
	listStackResourcesPages []*cloudformation.ListStackResourcesOutput
	listExportsPages        []*cloudformation.ListExportsOutput
	describeErr             error
	listStacksErr           error
	listResourcesErr        error
	listExportsErr          error
	stacksPage              int
	resourcePage            int
	exportPage              int
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
