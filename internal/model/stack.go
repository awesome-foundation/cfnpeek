package model

import "encoding/xml"

// StackInfo holds all inspectable data for a CloudFormation stack.
type StackInfo struct {
	XMLName   xml.Name   `json:"-" yaml:"-" toml:"-" xml:"stack"`
	StackName string     `json:"stack_name" yaml:"stack_name" toml:"stack_name" xml:"stack_name"`
	StackID   string     `json:"stack_id" yaml:"stack_id" toml:"stack_id" xml:"stack_id"`
	Status    string     `json:"status" yaml:"status" toml:"status" xml:"status"`
	Resources []Resource `json:"resources,omitempty" yaml:"resources,omitempty" toml:"resources,omitempty" xml:"resources>resource,omitempty"`
	Outputs   []Output   `json:"outputs,omitempty" yaml:"outputs,omitempty" toml:"outputs,omitempty" xml:"outputs>output,omitempty"`
	Exports   []Export   `json:"exports,omitempty" yaml:"exports,omitempty" toml:"exports,omitempty" xml:"exports>export,omitempty"`
}

type Resource struct {
	LogicalID   string `json:"logical_id" yaml:"logical_id" toml:"logical_id" xml:"logical_id"`
	PhysicalID  string `json:"physical_id" yaml:"physical_id" toml:"physical_id" xml:"physical_id"`
	Type        string `json:"type" yaml:"type" toml:"type" xml:"type"`
	Status      string `json:"status" yaml:"status" toml:"status" xml:"status"`
	LastUpdated string `json:"last_updated" yaml:"last_updated" toml:"last_updated" xml:"last_updated"`
}

type Output struct {
	Key         string `json:"key" yaml:"key" toml:"key" xml:"key"`
	Value       string `json:"value" yaml:"value" toml:"value" xml:"value"`
	Description string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty" xml:"description,omitempty"`
	ExportName  string `json:"export_name,omitempty" yaml:"export_name,omitempty" toml:"export_name,omitempty" xml:"export_name,omitempty"`
}

type Export struct {
	Name  string `json:"name" yaml:"name" toml:"name" xml:"name"`
	Value string `json:"value" yaml:"value" toml:"value" xml:"value"`
}

// StackSummary is a brief overview of a stack, used by the ls command.
type StackSummary struct {
	StackName   string `json:"stack_name" yaml:"stack_name" toml:"stack_name" xml:"stack_name"`
	StackID     string `json:"stack_id" yaml:"stack_id" toml:"stack_id" xml:"stack_id"`
	Status      string `json:"status" yaml:"status" toml:"status" xml:"status"`
	CreatedAt   string `json:"created_at" yaml:"created_at" toml:"created_at" xml:"created_at"`
	UpdatedAt   string `json:"updated_at,omitempty" yaml:"updated_at,omitempty" toml:"updated_at,omitempty" xml:"updated_at,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty" xml:"description,omitempty"`
}

// StackList wraps a list of stack summaries for formatted output.
type StackList struct {
	XMLName xml.Name       `json:"-" yaml:"-" toml:"-" xml:"stacks"`
	Stacks  []StackSummary `json:"stacks" yaml:"stacks" toml:"stacks" xml:"stack"`
}
