package filter_test

import (
	"testing"

	"github.com/awesome-foundation/cfnpeek/internal/filter"
	"github.com/awesome-foundation/cfnpeek/internal/model"
)

var sampleResources = []model.Resource{
	{LogicalID: "MyInstance", Type: "AWS::EC2::Instance", Status: "CREATE_COMPLETE"},
	{LogicalID: "MyBucket", Type: "AWS::S3::Bucket", Status: "CREATE_COMPLETE"},
	{LogicalID: "MySG", Type: "AWS::EC2::SecurityGroup", Status: "CREATE_COMPLETE"},
	{LogicalID: "MyVPC", Type: "AWS::EC2::VPC", Status: "CREATE_COMPLETE"},
}

var sampleOutputs = []model.Output{
	{Key: "VpcId", Value: "vpc-abc123"},
	{Key: "BucketArn", Value: "arn:aws:s3:::my-bucket"},
	{Key: "SomeOutput", Value: "value-with-vpc-ref"},
}

var sampleExports = []model.Export{
	{Name: "test-stack-VpcId", Value: "vpc-abc123"},
	{Name: "test-stack-BucketArn", Value: "arn:aws:s3:::my-bucket"},
}

func TestResources(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    []string
	}{
		{
			name:    "lowercase ec2 matches all EC2 resources",
			pattern: "ec2",
			want:    []string{"MyInstance", "MySG", "MyVPC"},
		},
		{
			name:    "uppercase EC2 matches case-insensitively",
			pattern: "EC2",
			want:    []string{"MyInstance", "MySG", "MyVPC"},
		},
		{
			name:    "s3 matches only S3 resources",
			pattern: "s3",
			want:    []string{"MyBucket"},
		},
		{
			name:    "instance filters to EC2::Instance only",
			pattern: "instance",
			want:    []string{"MyInstance"},
		},
		{
			name:    "no match returns nil slice",
			pattern: "lambda",
			want:    []string{},
		},
		{
			name:    "full type string matches exactly",
			pattern: "AWS::EC2::VPC",
			want:    []string{"MyVPC"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := filter.Resources(sampleResources, tc.pattern)
			gotIDs := logicalIDs(got)
			if len(gotIDs) != len(tc.want) {
				t.Fatalf("expected %d resources %v, got %d: %v", len(tc.want), tc.want, len(gotIDs), gotIDs)
			}
			for _, want := range tc.want {
				if !contains(gotIDs, want) {
					t.Errorf("expected logical ID %q in results, got %v", want, gotIDs)
				}
			}
		})
	}
}

func TestOutputs(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    []string
	}{
		{
			name:    "lowercase vpc matches key and value",
			pattern: "vpc",
			want:    []string{"VpcId", "SomeOutput"},
		},
		{
			name:    "uppercase VPC matches case-insensitively",
			pattern: "VPC",
			want:    []string{"VpcId", "SomeOutput"},
		},
		{
			name:    "bucket matches by key",
			pattern: "bucket",
			want:    []string{"BucketArn"},
		},
		{
			name:    "s3 matches by value",
			pattern: "s3",
			want:    []string{"BucketArn"},
		},
		{
			name:    "no match returns nil slice",
			pattern: "nonexistent",
			want:    []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := filter.Outputs(sampleOutputs, tc.pattern)
			gotKeys := outputKeys(got)
			if len(gotKeys) != len(tc.want) {
				t.Fatalf("expected %d outputs %v, got %d: %v", len(tc.want), tc.want, len(gotKeys), gotKeys)
			}
			for _, want := range tc.want {
				if !contains(gotKeys, want) {
					t.Errorf("expected output key %q in results, got %v", want, gotKeys)
				}
			}
		})
	}
}

func TestExports(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    []string
	}{
		{
			name:    "lowercase vpc matches export name",
			pattern: "vpc",
			want:    []string{"test-stack-VpcId"},
		},
		{
			name:    "uppercase VPC matches case-insensitively",
			pattern: "VPC",
			want:    []string{"test-stack-VpcId"},
		},
		{
			name:    "bucket matches by name",
			pattern: "bucket",
			want:    []string{"test-stack-BucketArn"},
		},
		{
			name:    "s3 matches by value",
			pattern: "s3",
			want:    []string{"test-stack-BucketArn"},
		},
		{
			name:    "no match returns nil slice",
			pattern: "nonexistent",
			want:    []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := filter.Exports(sampleExports, tc.pattern)
			gotNames := exportNames(got)
			if len(gotNames) != len(tc.want) {
				t.Fatalf("expected %d exports %v, got %d: %v", len(tc.want), tc.want, len(gotNames), gotNames)
			}
			for _, want := range tc.want {
				if !contains(gotNames, want) {
					t.Errorf("expected export name %q in results, got %v", want, gotNames)
				}
			}
		})
	}
}

func TestResourcesDoesNotMutateInput(t *testing.T) {
	original := make([]model.Resource, len(sampleResources))
	copy(original, sampleResources)

	filter.Resources(sampleResources, "ec2")

	for i, r := range sampleResources {
		if r != original[i] {
			t.Errorf("input slice was mutated at index %d", i)
		}
	}
}

// helpers

func logicalIDs(resources []model.Resource) []string {
	ids := make([]string, len(resources))
	for i, r := range resources {
		ids[i] = r.LogicalID
	}
	return ids
}

func outputKeys(outputs []model.Output) []string {
	keys := make([]string, len(outputs))
	for i, o := range outputs {
		keys[i] = o.Key
	}
	return keys
}

func exportNames(exports []model.Export) []string {
	names := make([]string, len(exports))
	for i, e := range exports {
		names[i] = e.Name
	}
	return names
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
