package jsonplan

import (
	"github.com/hashicorp/terraform/configs/configload"
	"github.com/hashicorp/terraform/plans"
	"github.com/hashicorp/terraform/states"
)

// FormatVersion represents the version of the json format and will be incremented
// for any change to this format that requires changes to a consuming parser
const FormatVersion = "0.1"

// Plan is the top-level representation of the json format of a plan
// It includes the complete config and current state
type Plan struct {
	FormatVersion   string            `json:"format_version"`
	PriorState      Values            `json:"prior_state,omitempty"`
	Config          map[string]Config `json:"configuration"`
	PlannedValues   Values            `json:"planned_values"`
	ProposedUnknown Values            `json:"proposed_unknown"`
	ResourceChanges []ResourceChange  `json:"resource_changes"`
	OutputChanges   map[string]Change `json:"output_changes"`
}

// Change is the representation of a proposed change for an object
type Change struct {
	// Actions are the actions that will be taken on the object selected by
	// the properties below.
	// Valid actions values are:
	//    ["no-op"]
	//    ["create"]
	//    ["read"]
	//    ["update"]
	//    ["delete", "create"]
	//    ["create", "delete"]
	//    ["delete"]
	// The two "replace" actions are represented in this way to allow callers to
	// e.g. just scan the list for "delete" to recognize all three situations
	// where the object will be deleted, allowing for any new deletion
	// combinations that might be added in future.
	Actions []string

	// Before and After are representations of the object value both before and
	// after the action. For ["create"] and ["delete"] actions, either "before"
	// or "after" is unset (respectively). For ["no-op"], the before and after
	// values are identical. The "after" value will be incomplete if there are
	// values within it that won't be known until after apply.
	Before Value
	After  Value
}

// Values is the common representation of resolved values for both the prior
// state (which is always complete) and the planned new state
type Values struct {
	Outputs      map[string]Output
	RootModule   Module
	ChildModules []Module
}

// Value .... is where I got lost.
type Value struct{}

// Config represents the complete configuration source
type Config struct {
	ProviderConfigs []ProviderConfig `json:"provider_config"`
	RootModule      ConfigRootModule `json:"root_module"`
}

// ProviderConfig describes all of the provider configurations throughout the
// configuration tree, flattened into a single map for convenience since
// provider configurations are the one concept in Terraform that can span across
// module boundaries.
type ProviderConfig struct {
	Name          string
	Alias         string
	ModuleAddress string
	Expressions   Expressions
}

type ConfigRootModule struct {
	Outputs     []map[string]Output
	Resources   []Resource
	ModuleCalls []ModuleCall
}

// Resource is the representation of a resource in the json plan
type Resource struct {
	// Address is the absolute resource address
	Address string `json:"address"`

	// "managed" or "data"
	Mode string `json:"mode"`

	Type string `json:"type"`
	Name string `json:"name"`

	// Index is omitted for a resource not using `count` or `for_each`
	Index int `json:"index,omitempty"`

	// This is only the provider name, not a provider configuration address.
	//
	// It is included to allow the property "type" to be interpreted
	// unambiguously in the unusual situation where a provider offers a resource
	// type whose name does not start with its own name, such as the "googlebeta"
	// provider offering "google_compute_instance".
	ProviderName string `json:"provider_name"`

	// "schema_version" indicates which version of the resource type schema the
	// "values" property conforms to.
	SchemaVersion int `json:"schema_version"`

	// "values" is the JSON representation of the attribute values of the
	// resource, whose structure depends on the resource type schema.
	// Any unknown values are omitted or set to null, making them indistinguishable
	// from absent values.
	Values []byte `json:"values"`
}

//
type ResourceChange struct {
	// Address is the absolute resource address
	Address string `json:"address,omitempty"`

	// ModuleAddress is the module portion of the above address. Omitted if the
	// instance is in the root module.
	ModuleAddress string `json:"module_address,omitempty"`

	// "managed" or "data"
	Mode string

	Type  string
	Name  string
	Index string

	// "deposed", if set, indicates that this action applies to a "deposed"
	// object of the given instance rather than to its "current" object. Omitted
	// for changes to the current object.
	Deposed bool `json:"deposed,omitempty"`

	// Change describes the change that will be made to this object
	Change Change
}

// Module is the representation of a module in state
// This can be the root module or a child module
type Module struct {
	Resources []Resource

	// Address is the absolute module address, omitted for the root module
	Address string `json:"address,omitempty"`

	// Each module object can optionally have its own nested "child_modules",
	// recursively describing the full module tree.
	ChildModules []Module `json:"child_modules,omitempty`
}

type ModuleCall struct {
	ResolvedSource    string      `json:"resolved_source"`
	Expressions       Expressions `json:"expressions"`
	CountExpression   Expression  `json:"count_expression"`
	ForEachExpression Expression  `json:"for_each_expression"`
	Module            Module      `json:""`
}

type Output struct {
	Sensitive bool
	Value     string
}

type ConfigOutput struct {
	Sensitive  bool
	Expression Expression
}

type Expression struct {
	ConstantValue string   `json:"constant_value,omitempty"`
	References    []string `json:"references,omitempty`
	Source        Source   `json:source`
}

type Expressions struct {
	Expression map[string]Expression
}

type Source struct {
	FileName string `json:"filename"`
	Start    string `json:"start"`
	End      string `json:"end"`
}

// Marshall returns the json encoding of a terraform plan
func Marshall(c *configload.Snapshot, p *plans.Plan, s *states.State) ([]byte, error) {
	return nil, nil
}
