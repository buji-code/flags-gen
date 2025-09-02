// Package parser provides functionality for parsing Go source files and extracting
// struct information for flag generation. It uses Go's AST package to analyze
// struct definitions marked with +flags-gen annotations and extract field metadata
// including types, tags, and documentation comments.
package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strconv"
	"strings"

	"github.com/yuvalwz/flags-gen/pkg/types"
)

// Parser handles parsing Go source files for structs with flags-gen annotations
type Parser struct {
	fileSet *token.FileSet
}

// New creates a new Parser instance
func New() *Parser {
	return &Parser{
		fileSet: token.NewFileSet(),
	}
}

// ParseFile parses a Go source file and returns structs marked with +flags-gen
func (p *Parser) ParseFile(filename string) ([]types.StructInfo, error) {
	src, err := parser.ParseFile(p.fileSet, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filename, err)
	}

	var structs []types.StructInfo
	
	// Walk through all declarations in the file
	for _, decl := range src.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						// Check if this struct has the +flags-gen annotation
						if p.hasAnnotation(genDecl.Doc) {
							structInfo, err := p.parseStruct(typeSpec.Name.Name, structType, src.Name.Name)
							if err != nil {
								return nil, fmt.Errorf("failed to parse struct %s: %w", typeSpec.Name.Name, err)
							}
							structs = append(structs, structInfo)
						}
					}
				}
			}
		}
	}

	return structs, nil
}

// hasAnnotation checks if the comment group contains +flags-gen annotation
func (p *Parser) hasAnnotation(commentGroup *ast.CommentGroup) bool {
	if commentGroup == nil {
		return false
	}
	
	for _, comment := range commentGroup.List {
		if strings.Contains(comment.Text, "+flags-gen") {
			return true
		}
	}
	return false
}

// parseStruct parses a struct and extracts field information for flag generation
func (p *Parser) parseStruct(name string, structType *ast.StructType, packageName string) (types.StructInfo, error) {
	structInfo := types.StructInfo{
		Name:        name,
		PackageName: packageName,
		Fields:      make([]types.FieldInfo, 0),
		Imports:     make([]string, 0),
	}

	imports := make(map[string]bool)

	for _, field := range structType.Fields.List {
		// Skip embedded fields or fields without names
		if len(field.Names) == 0 {
			continue
		}

		for _, fieldName := range field.Names {
			// Skip unexported fields
			if !ast.IsExported(fieldName.Name) {
				continue
			}

			fieldInfo, err := p.parseField(fieldName.Name, field)
			if err != nil {
				return structInfo, fmt.Errorf("failed to parse field %s: %w", fieldName.Name, err)
			}

			// Add required imports based on field type
			if fieldInfo.Type == "time.Duration" {
				imports["time"] = true
			}

			// Set flag method and default value code
			if method, exists := types.GetFlagMethod(fieldInfo.Type); exists {
				fieldInfo.FlagMethod = method
				fieldInfo.DefaultValueCode = p.formatDefaultValueCode(fieldInfo.DefaultValue, fieldInfo.Type)
			}

			structInfo.Fields = append(structInfo.Fields, fieldInfo)
		}
	}

	// Convert imports map to slice
	for imp := range imports {
		structInfo.Imports = append(structInfo.Imports, imp)
	}

	return structInfo, nil
}

// parseField extracts information from a single struct field
func (p *Parser) parseField(name string, field *ast.Field) (types.FieldInfo, error) {
	fieldInfo := types.FieldInfo{
		Name: name,
	}

	// Parse field type
	fieldType, err := p.parseType(field.Type)
	if err != nil {
		return fieldInfo, fmt.Errorf("failed to parse type for field %s: %w", name, err)
	}
	fieldInfo.Type = fieldType

	// Parse struct tags
	if field.Tag != nil {
		tag := strings.Trim(field.Tag.Value, "`")
		fieldInfo.JSONTag = p.extractJSONTag(tag)
		fieldInfo.FlagName = p.deriveFlagName(name, fieldInfo.JSONTag)
		
		// Look for default values in tags
		fieldInfo.DefaultValue = p.extractDefaultFromTag(tag, fieldType)
	} else {
		fieldInfo.FlagName = p.deriveFlagName(name, "")
	}

	// Parse field comments for description
	fieldInfo.Description = p.parseFieldComment(field.Comment, field.Doc)

	return fieldInfo, nil
}

// parseType converts an ast.Expr representing a type to a string
func (p *Parser) parseType(expr ast.Expr) (string, error) {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name, nil
	case *ast.SelectorExpr:
		pkg, err := p.parseType(t.X)
		if err != nil {
			return "", err
		}
		return pkg + "." + t.Sel.Name, nil
	case *ast.ArrayType:
		elemType, err := p.parseType(t.Elt)
		if err != nil {
			return "", err
		}
		return "[]" + elemType, nil
	default:
		return "", fmt.Errorf("unsupported type: %T", expr)
	}
}

// extractJSONTag extracts the json tag value from struct tag
func (p *Parser) extractJSONTag(tag string) string {
	re := regexp.MustCompile(`json:"([^"]*)"`)
	matches := re.FindStringSubmatch(tag)
	if len(matches) > 1 {
		// Remove omitempty and other options
		parts := strings.Split(matches[1], ",")
		if parts[0] != "" {
			return parts[0]
		}
	}
	return ""
}

// extractDefaultFromTag extracts default values from struct tags
func (p *Parser) extractDefaultFromTag(tag, fieldType string) interface{} {
	re := regexp.MustCompile(`default:"([^"]*)"`)
	matches := re.FindStringSubmatch(tag)
	if len(matches) > 1 {
		defaultStr := matches[1]
		return p.parseDefaultValue(defaultStr, fieldType)
	}
	return nil
}

// parseDefaultValue converts string default value to appropriate type
func (p *Parser) parseDefaultValue(value, fieldType string) interface{} {
	switch fieldType {
	case "string":
		return value
	case "int", "int32", "int64":
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	case "bool":
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	case "[]string":
		if value != "" {
			return strings.Split(value, ",")
		}
		return []string{}
	case "time.Duration":
		return value // Keep as string, will be parsed later
	}
	return value
}

// deriveFlagName creates a flag name from field name and json tag
func (p *Parser) deriveFlagName(fieldName, jsonTag string) string {
	if jsonTag != "" {
		return p.toKebabCase(jsonTag)
	}
	return p.toKebabCase(fieldName)
}

// toKebabCase converts camelCase to kebab-case
func (p *Parser) toKebabCase(s string) string {
	// Handle sequences of capital letters (e.g., HTTPPort -> HTTP-Port)
	re1 := regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	result := re1.ReplaceAllString(s, "${1}-${2}")
	
	// Handle normal camelCase (e.g., camelCase -> camel-Case)
	re2 := regexp.MustCompile(`([a-z])([A-Z])`)
	result = re2.ReplaceAllString(result, "${1}-${2}")
	
	return strings.ToLower(result)
}

// parseFieldComment extracts description from field comments
func (p *Parser) parseFieldComment(comment *ast.CommentGroup, doc *ast.CommentGroup) string {
	var description string
	
	// Check doc comment first (appears before the field)
	if doc != nil {
		for _, c := range doc.List {
			text := strings.TrimPrefix(c.Text, "//")
			text = strings.TrimPrefix(text, "/*")
			text = strings.TrimSuffix(text, "*/")
			text = strings.TrimSpace(text)
			
			// Skip annotations like +optional
			if strings.HasPrefix(text, "+") {
				continue
			}
			
			if description == "" {
				description = text
			} else {
				description += " " + text
			}
		}
	}
	
	// Check inline comment if no doc comment found
	if description == "" && comment != nil {
		for _, c := range comment.List {
			text := strings.TrimPrefix(c.Text, "//")
			text = strings.TrimSpace(text)
			description = text
			break // Only take the first inline comment
		}
	}
	
	return description
}

// formatDefaultValueCode formats a default value for code generation
func (p *Parser) formatDefaultValueCode(value interface{}, fieldType string) string {
	if value == nil {
		return p.getZeroValue(fieldType)
	}
	
	switch fieldType {
	case "string":
		return fmt.Sprintf(`"%s"`, value)
	case "[]string":
		if slice, ok := value.([]string); ok {
			quoted := make([]string, len(slice))
			for i, s := range slice {
				quoted[i] = fmt.Sprintf(`"%s"`, s)
			}
			return fmt.Sprintf("[]string{%s}", strings.Join(quoted, ", "))
		}
		return `[]string{}`
	case "time.Duration":
		if str, ok := value.(string); ok {
			return fmt.Sprintf("%s*time.Second", strings.TrimSuffix(str, "s"))
		}
		return "0"
	default:
		return fmt.Sprintf("%v", value)
	}
}

// getZeroValue returns the zero value for a given type
func (p *Parser) getZeroValue(fieldType string) string {
	switch fieldType {
	case "string":
		return `""`
	case "int", "int32", "int64", "uint", "uint32", "uint64", "float32", "float64":
		return "0"
	case "bool":
		return "false"
	case "[]string", "[]int":
		return fmt.Sprintf("%s{}", fieldType)
	case "time.Duration":
		return "0"
	default:
		return `""`
	}
}