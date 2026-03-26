package formatter_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"strings"
	"testing"

	"github.com/awesome-foundation/cfnpeek/internal/formatter"
	"github.com/awesome-foundation/cfnpeek/internal/model"
)

func testData() *model.StackInfo {
	return &model.StackInfo{
		StackName: "test-stack",
		StackID:   "arn:aws:cloudformation:us-east-1:123456789:stack/test-stack/guid",
		Status:    "CREATE_COMPLETE",
		Resources: []model.Resource{
			{
				LogicalID:   "MyBucket",
				PhysicalID:  "test-stack-mybucket-abc123",
				Type:        "AWS::S3::Bucket",
				Status:      "CREATE_COMPLETE",
				LastUpdated: "2026-01-01T00:00:00Z",
			},
		},
		Outputs: []model.Output{
			{
				Key:        "BucketArn",
				Value:      "arn:aws:s3:::test-stack-mybucket-abc123",
				ExportName: "test-stack-BucketArn",
			},
		},
		Exports: []model.Export{
			{
				Name:  "test-stack-BucketArn",
				Value: "arn:aws:s3:::test-stack-mybucket-abc123",
			},
		},
	}
}

func TestGetValidFormat(t *testing.T) {
	for _, name := range []string{"json", "yaml", "toml", "xml", "ini", "csv", "table"} {
		f, err := formatter.Get(name)
		if err != nil {
			t.Errorf("Get(%q) returned error: %v", name, err)
		}
		if f == nil {
			t.Errorf("Get(%q) returned nil formatter", name)
		}
	}
}

func TestGetInvalidFormat(t *testing.T) {
	_, err := formatter.Get("invalid")
	if err == nil {
		t.Error("Get(\"invalid\") should return error")
	}
}

func TestJSONRoundtrip(t *testing.T) {
	data := testData()
	var buf bytes.Buffer
	f, _ := formatter.Get("json")
	if err := f.Format(&buf, data); err != nil {
		t.Fatal(err)
	}

	var decoded model.StackInfo
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("JSON output is not valid: %v\nOutput:\n%s", err, buf.String())
	}
	if decoded.StackName != "test-stack" {
		t.Errorf("expected stack name 'test-stack', got %q", decoded.StackName)
	}
	if len(decoded.Resources) != 1 {
		t.Errorf("expected 1 resource, got %d", len(decoded.Resources))
	}
}

func TestYAMLContainsFields(t *testing.T) {
	data := testData()
	var buf bytes.Buffer
	f, _ := formatter.Get("yaml")
	if err := f.Format(&buf, data); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"stack_name:", "test-stack", "MyBucket", "AWS::S3::Bucket"} {
		if !strings.Contains(out, want) {
			t.Errorf("YAML output missing %q", want)
		}
	}
}

func TestXMLValid(t *testing.T) {
	data := testData()
	var buf bytes.Buffer
	f, _ := formatter.Get("xml")
	if err := f.Format(&buf, data); err != nil {
		t.Fatal(err)
	}

	var decoded model.StackInfo
	if err := xml.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("XML output is not valid: %v\nOutput:\n%s", err, buf.String())
	}
}

func TestTOMLContainsFields(t *testing.T) {
	data := testData()
	var buf bytes.Buffer
	f, _ := formatter.Get("toml")
	if err := f.Format(&buf, data); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"name", "test-stack", "MyBucket"} {
		if !strings.Contains(out, want) {
			t.Errorf("TOML output missing %q", want)
		}
	}
}

func TestINIContainsSections(t *testing.T) {
	data := testData()
	var buf bytes.Buffer
	f, _ := formatter.Get("ini")
	if err := f.Format(&buf, data); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"[stack]", "[resource.MyBucket]", "[outputs]", "[exports]"} {
		if !strings.Contains(out, want) {
			t.Errorf("INI output missing section %q", want)
		}
	}
}

func TestCSVContainsHeaders(t *testing.T) {
	data := testData()
	var buf bytes.Buffer
	f, _ := formatter.Get("csv")
	if err := f.Format(&buf, data); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"logical_id", "physical_id", "key", "value"} {
		if !strings.Contains(out, want) {
			t.Errorf("CSV output missing header %q", want)
		}
	}
}

func TestTableContainsLabels(t *testing.T) {
	data := testData()
	var buf bytes.Buffer
	f, _ := formatter.Get("table")
	if err := f.Format(&buf, data); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"Stack: test-stack", "Resources (1)", "Outputs (1)", "Exports (1)"} {
		if !strings.Contains(out, want) {
			t.Errorf("table output missing %q", want)
		}
	}
}

func TestEmptySections(t *testing.T) {
	data := &model.StackInfo{
		StackName: "empty-stack",
		StackID:   "arn:aws:cloudformation:us-east-1:123456789:stack/empty-stack/guid",
		Status:    "CREATE_COMPLETE",
	}

	for _, name := range []string{"json", "yaml", "toml", "xml", "ini", "csv", "table"} {
		var buf bytes.Buffer
		f, _ := formatter.Get(name)
		if err := f.Format(&buf, data); err != nil {
			t.Errorf("%s formatter failed on empty data: %v", name, err)
		}
	}
}
