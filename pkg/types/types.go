// Package types defines the core data structures and utilities used by the
// flags-gen tool for representing parsed struct information and supported
// flag type mappings.
package types

const (
	// Type constants.
	TypeString       = "string"
	TypeBool         = "bool"
	TypeInt          = "int"
	TypeInt32        = "int32"
	TypeInt64        = "int64"
	TypeStringSlice  = "[]string"
	TypeTimeDuration = "time.Duration"
)

// FieldInfo represents information about a struct field that needs flag generation.
type FieldInfo struct {
	Name             string
	Type             string
	JSONTag          string
	FlagName         string
	Description      string
	DefaultValue     interface{}
	DefaultValueCode string
	Required         bool
	ShortFlag        string
	FlagMethod       string
}

// StructInfo represents information about a struct that needs flag generation.
type StructInfo struct {
	Name        string
	PackageName string
	Fields      []FieldInfo
	Imports     []string
}

// SupportedTypes maps Go types to their pflags method names.
var SupportedTypes = map[string]string{
	"string":        "StringVar",
	"int":           "IntVar",
	"int32":         "Int32Var",
	"int64":         "Int64Var",
	"uint":          "UintVar",
	"uint32":        "Uint32Var",
	"uint64":        "Uint64Var",
	"bool":          "BoolVar",
	"float32":       "Float32Var",
	"float64":       "Float64Var",
	"[]string":      "StringSliceVar",
	"[]int":         "IntSliceVar",
	"time.Duration": "DurationVar",
}

// GetFlagMethod returns the appropriate pflags method for a given type.
func GetFlagMethod(fieldType string) (string, bool) {
	method, exists := SupportedTypes[fieldType]
	return method, exists
}

// HasShortFlag returns true if the field supports short flags (single character flags).
func HasShortFlag(fieldType string) bool {
	// Only simple types typically get short flags to avoid confusion
	switch fieldType {
	case "string", "int", "bool":
		return true
	default:
		return false
	}
}
