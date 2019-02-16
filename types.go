package main

import (
	"time"
)

// Indicator that is used in alerts, pipelines, etc.
type Indicator struct {
	Id          string `json:"id,omitempty"`
	Type        string `json:"type,omitempty"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"`
	Author      string `json:"author,omitempty"`
	Source      string `json:"source,omitempty"`

	// An indicator would never have probability 0.0, so zero means not
	// specified.
	Probability float32 `json:"probability,omitempty"`
}

// Indicators is a collection of flat indicators. This is now
// "IndicatorDefinitions", but we are probably using it someplace...
type Indicators struct {
	Description string       `json:"description,omitempty"`
	Version     string       `json:"version,omitempty"`
	Indicators  []*Indicator `json:"indicators,omitempty"`
}

// TypeValuePair is a type/value pair. Duh.
type TypeValuePair struct {
	Type  string `json:"type" required:"true"`
	Value string `json:"value" required:"true"`
}

// IOCDefinition is used as the branches and leaves of the tree
type IOCDefinition struct {
	Type  string `json:"type, omitempty"`
	Value string `json:"value, omitempty"`
}

// IndicatorDefinition contains the data comprising an IOCDefinition
type IndicatorDefinition struct {
	ID            string         `json:"id" required:"true"`
	AuthoredDate  time.Time      `json:"authored_date" required:"true"`
	LastModified  time.Time      `json:"last_modified" required:"true"`
	Description   string         `json:"description" required:"true"`
	Source        string         `json:"source" required:"true"`
	Author        string         `json:"author" required:"true"`
	Category      string         `json:"category" required:"true"`
	IOCDefinition *IOCDefinition `json:"ioc_definition" required:"true"`
}

// Pattern is the pattern to match on
// The Type is the type of event property to match, e.g. "country"
// Value is the value to match
// Value2 is a second value to match, e.g. required for a range match
// Match is the type of match to perform:
//    - string (string match of Value, the default if Match is not specified)
//    - int (an integer match of Value)
//    - range (an integer range match of Value-Value2 inclusive)
//    - dns (a DNS hostname match of Value)
type Pattern struct {
	Type   string `json:"type,omitempty"`
	Value  string `json:"value,omitempty"`
	Value2 string `json:"value2,omitempty"`
	Match  string `json:"match,omitempty"`
}

// IndicatorNode is a node in a boolean tree.
// A node may have children, in which case it must have an Operator, or
//  it might be a leaf node, in which case it must have a Pattern to match on.
// A node may be just a reference to another 'concrete' node - you cannot
//  reference a reference node though (there is no point)
// Children are specified in the IOCs definition file(s); links to Parents are
//  created at IOC def load time.
// This struct is used for both the IOC def file(s) and the runtime lookups.
type IndicatorNode struct {
	ID          string           `json:"id,omitempty"`
	Comment     string           `json:"comment,omitempty"`
	Ref         string           `json:"ref,omitempty"`
	Operator    string           `json:"operator,omitempty"` // OR|AND|NOT
	Indicator   *Indicator       `json:"indicator,omitempty"`
	Parents     []*IndicatorNode `json:"parents,omitempty"`
	Children    []*IndicatorNode `json:"children,omitempty"`
	SiblingNots []int            `json:"siblingnots,omitempty"`
	Pattern     *Pattern         `json:"pattern,omitempty"`
}

// IndicatorDefinitions defines the file format of IOC definitions.
// IOCs could be defined with multiple of such files.
type IndicatorDefinitions struct {
	Description string           `json:"description,omitempty"`
	Version     string           `json:"version,omitempty"`
	Definitions []*IndicatorNode `json:"definitions,omitempty"`
}
