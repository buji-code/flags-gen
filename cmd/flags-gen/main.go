package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yuvalwz/flags-gen/pkg/generator"
	"github.com/yuvalwz/flags-gen/pkg/parser"
)

var (
	inputFile  string
	outputFile string
	version    = "dev"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "flags-gen",
		Short: "Generate pflags AddFlags methods from Go structs",
		Long: `flags-gen is a code generation tool that parses Go structs marked with +flags-gen
annotations and generates corresponding pflags AddFlags methods for CLI applications.

Example:
  flags-gen -i types.go -o flags_gen.go
  flags-gen --input=./pkg/types/config.go --output=./pkg/types/flags.go`,
		RunE: runFlagsGen,
	}

	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input Go file containing structs with +flags-gen annotations (required)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for generated flags code (optional, defaults to <input>_flags.go)")
	rootCmd.MarkFlagRequired("input")

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("flags-gen version %s\n", version)
		},
	}

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runFlagsGen(cmd *cobra.Command, args []string) error {
	if inputFile == "" {
		return fmt.Errorf("input file is required")
	}

	// Validate and clean input file path
	cleanInputFile, err := validateFilePath(inputFile)
	if err != nil {
		return fmt.Errorf("invalid input file path: %w", err)
	}
	inputFile = cleanInputFile

	// Validate input file exists and is accessible
	fileInfo, err := os.Stat(inputFile)
	if os.IsNotExist(err) {
		return fmt.Errorf("input file %s does not exist\n\nTip: Make sure the file path is correct and the file has a .go extension", inputFile)
	}
	if err != nil {
		return fmt.Errorf("cannot access input file %s: %w", inputFile, err)
	}

	// Check file size to prevent DoS
	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if fileInfo.Size() > maxFileSize {
		return fmt.Errorf("input file %s is too large (%d bytes), maximum allowed size is %d bytes", inputFile, fileInfo.Size(), maxFileSize)
	}

	// Ensure input file has .go extension
	if !strings.HasSuffix(strings.ToLower(inputFile), ".go") {
		return fmt.Errorf("input file must be a Go source file (.go extension)")
	}

	// Generate output file name if not provided
	if outputFile == "" {
		dir := filepath.Dir(inputFile)
		base := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
		outputFile = filepath.Join(dir, base+"_flags.go")
	}

	// Validate output file path
	cleanOutputFile, err := validateFilePath(outputFile)
	if err != nil {
		return fmt.Errorf("invalid output file path: %w", err)
	}
	outputFile = cleanOutputFile

	// Parse the input file
	p := parser.New()
	structs, err := p.ParseFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to parse input file: %w", err)
	}

	if len(structs) == 0 {
		return fmt.Errorf("no structs with +flags-gen annotation found in %s", inputFile)
	}

	// Generate flags code for all structs
	g := generator.New()
	var allGenerated []string

	for _, structInfo := range structs {
		generated, err := g.GenerateFlags(structInfo)
		if err != nil {
			return fmt.Errorf("failed to generate flags for struct %s: %w", structInfo.Name, err)
		}
		allGenerated = append(allGenerated, generated)
	}

	// Write output file
	output := strings.Join(allGenerated, "\n\n")
	if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Generated flags code for %d struct(s) in %s\n", len(structs), outputFile)
	return nil
}

// validateFilePath validates and cleans a file path to prevent path traversal attacks
func validateFilePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("file path cannot be empty")
	}

	// Clean the path to resolve any relative components
	cleanPath := filepath.Clean(path)

	// Check for suspicious patterns that might indicate path traversal
	if strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("path contains directory traversal patterns")
	}

	// Convert to absolute path if relative
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	return absPath, nil
}