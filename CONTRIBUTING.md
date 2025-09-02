# Contributing to flags-gen

Thank you for your interest in contributing to flags-gen! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

This project adheres to a code of conduct that promotes respectful collaboration. Please read and follow our code of conduct in all interactions.

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.20+**: [Download and install Go](https://golang.org/dl/)
- **Git**: [Install Git](https://git-scm.com/downloads)
- **Make**: Most systems have this installed, or you can [install GNU Make](https://www.gnu.org/software/make/)

### Optional Tools

These tools are recommended for the best development experience:

- **golangci-lint**: [Install golangci-lint](https://golangci-lint.run/usage/install/) for comprehensive code linting
- **gosec**: `go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest` for security analysis
- **govulncheck**: `go install golang.org/x/vuln/cmd/govulncheck@latest` for vulnerability scanning

## Development Environment Setup

### 1. Fork and Clone

```bash
# Fork the repository on GitHub, then clone your fork
git clone https://github.com/your-username/flags-gen.git
cd flags-gen

# Add the original repository as upstream
git remote add upstream https://github.com/yuvalwz/flags-gen.git
```

### 2. Install Dependencies

```bash
# Download and tidy dependencies
make deps
```

### 3. Verify Setup

```bash
# Build the project
make build

# Run tests
make test

# Run all quality checks
make check
```

### 4. Development Workflow

```bash
# Create a feature branch
git checkout -b feature/your-feature-name

# Make your changes...

# Run quality checks before committing
make fmt      # Format code
make vet      # Run go vet
make lint     # Run linter (if installed)
make test     # Run tests

# Commit your changes
git add .
git commit -m "feat: add your feature description"

# Push to your fork
git push origin feature/your-feature-name
```

## Project Structure

Understanding the codebase structure will help you contribute effectively:

```
flags-gen/
├── cmd/flags-gen/          # CLI application entry point
│   ├── main.go            # Main command implementation
│   └── main_test.go       # CLI tests
├── pkg/                   # Public packages
│   ├── parser/            # Go AST parsing logic
│   │   ├── parser.go      # Main parser implementation
│   │   └── parser_test.go # Parser tests
│   ├── generator/         # Code generation logic
│   │   ├── generator.go   # Template-based code generator
│   │   └── generator_test.go
│   └── types/             # Type definitions and utilities
│       └── types.go       # Shared types and constants
├── internal/              # Private packages
│   └── testdata/          # Test fixtures and examples
│       ├── example.go     # Example input struct
│       └── example_flags.go # Generated output
├── docs/                  # Documentation (you'll create this)
├── Makefile              # Build and development tasks
├── go.mod                # Go module definition
├── .golangci.yml         # Linter configuration
└── README.md             # Project documentation
```

### Key Components

1. **Parser (`pkg/parser/`)**: Analyzes Go source files using the `go/ast` package to find structs with `+flags-gen` annotations
2. **Generator (`pkg/generator/`)**: Uses Go templates to generate `AddFlags` methods from parsed struct information
3. **Types (`pkg/types/`)**: Defines data structures and type mappings used throughout the application
4. **CLI (`cmd/flags-gen/`)**: Command-line interface using Cobra

## Development Guidelines

### Code Style

We follow standard Go conventions:

- **Formatting**: Use `gofmt` (run `make fmt`)
- **Naming**: Follow Go naming conventions (exported vs unexported, camelCase, etc.)
- **Comments**: Public functions and types must have comments starting with the name
- **Error Handling**: Always handle errors explicitly, wrap with context when appropriate

### Code Quality Standards

All code must pass these quality checks:

```bash
make fmt      # Code formatting
make vet      # Go vet static analysis
make lint     # golangci-lint (if available)
make test     # All tests must pass
make coverage # Maintain >80% test coverage
```

### Writing Tests

We use table-driven tests where appropriate:

```go
func TestParseType(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "simple string type",
            input:    "string",
            expected: "string",
            wantErr:  false,
        },
        // Add more test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := parseType(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("parseType() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if result != tt.expected {
                t.Errorf("parseType() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Documentation Standards

- **Public APIs**: Must have comprehensive Go doc comments
- **Complex Logic**: Include inline comments explaining non-obvious code
- **Examples**: Provide usage examples for new features
- **README Updates**: Update README.md for user-facing changes

### Git Commit Guidelines

We follow [Conventional Commits](https://conventionalcommits.org/) specification:

```bash
# Format
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]

# Examples
feat: add support for float32 and float64 types
fix: handle embedded structs correctly
docs: update installation instructions
test: add integration tests for CLI
refactor: simplify type mapping logic
chore: update dependencies to latest versions
```

#### Commit Types

- **feat**: New feature
- **fix**: Bug fix  
- **docs**: Documentation only changes
- **style**: Code style changes (formatting, etc.)
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **test**: Adding missing tests or correcting existing tests
- **chore**: Changes to build process or auxiliary tools
- **perf**: Performance improvements
- **ci**: CI/CD changes

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with race detection
make test-race

# Generate coverage report
make coverage

# Run benchmarks
make bench
```

### Test Categories

1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test component interactions
3. **CLI Tests**: Test command-line interface behavior
4. **Golden File Tests**: Compare generated output with expected files

### Writing New Tests

When adding new features:

1. Write tests for all new public functions
2. Include both success and error cases
3. Test edge cases and boundary conditions
4. Add integration tests for end-to-end workflows
5. Update existing tests if behavior changes

### Test Data

- Place test fixtures in `internal/testdata/`
- Use meaningful names for test files
- Include both simple and complex examples
- Update golden files when output format changes

## Making Changes

### Adding New Features

1. **Design First**: Open an issue to discuss the feature
2. **Start Small**: Break large features into smaller, reviewable changes
3. **Write Tests**: Add comprehensive test coverage
4. **Update Docs**: Update README and other documentation
5. **Backward Compatibility**: Avoid breaking existing functionality

### Bug Fixes

1. **Reproduce**: Create a test that demonstrates the bug
2. **Fix**: Implement the minimal fix required
3. **Verify**: Ensure the test passes and no regressions occur
4. **Document**: Update relevant documentation if needed

### Common Contribution Areas

#### 1. New Type Support

To add support for a new Go type:

1. Update `SupportedTypes` map in `pkg/types/types.go`
2. Add parsing logic in `pkg/parser/parser.go` if needed
3. Add generation logic in `pkg/generator/generator.go` if needed
4. Add comprehensive tests
5. Update documentation

#### 2. Parser Improvements

When improving the parser:

1. Add test cases in `pkg/parser/parser_test.go`
2. Update the parser logic in `pkg/parser/parser.go`
3. Test with complex struct examples
4. Ensure backward compatibility

#### 3. Generator Enhancements

When improving code generation:

1. Update the template in `pkg/generator/generator.go`
2. Add test cases for new generation patterns
3. Verify generated code compiles and works correctly
4. Test with various flag combinations

## Pull Request Process

### Before Submitting

1. **Update from upstream**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run all checks**:
   ```bash
   make check
   ```

3. **Test thoroughly**:
   ```bash
   make test
   make coverage
   ```

### Pull Request Checklist

- [ ] All tests pass
- [ ] Code coverage maintained or improved
- [ ] Documentation updated (README, code comments, etc.)
- [ ] Backward compatibility maintained
- [ ] Commit messages follow conventional format
- [ ] Changes are focused and atomic
- [ ] Self-review completed

### PR Description Template

```markdown
## Description
Brief description of changes and motivation.

## Type of Change
- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to change)
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Tests pass locally
- [ ] Documentation updated
```

### Review Process

1. **Automated Checks**: CI/CD pipeline runs tests and quality checks
2. **Code Review**: Maintainers review code for quality and design
3. **Feedback**: Address any requested changes
4. **Approval**: PR approved by maintainers
5. **Merge**: PR merged into main branch

## Release Process

Releases are handled by project maintainers:

1. **Version Bump**: Update version in relevant files
2. **Changelog**: Update CHANGELOG.md with release notes
3. **Tag**: Create and push version tag
4. **Release**: GitHub Actions builds and publishes release
5. **Announce**: Update documentation and announce release

## Getting Help

### Communication Channels

- **GitHub Issues**: For bugs, feature requests, and discussions
- **GitHub Discussions**: For questions and community discussion
- **Pull Requests**: For code review and collaboration

### Asking Questions

When asking for help:

1. **Search First**: Check existing issues and documentation
2. **Provide Context**: Include relevant code, error messages, and environment details
3. **Be Specific**: Clear, focused questions get better answers
4. **Be Patient**: Maintainers are volunteers with limited time

### Reporting Issues

When reporting bugs:

1. **Use the Issue Template**: Fill out all relevant sections
2. **Include Examples**: Provide minimal reproducible examples
3. **Environment Details**: Include Go version, OS, and flags-gen version
4. **Expected vs Actual**: Clearly describe what you expected vs what happened

## Recognition

Contributors are recognized in several ways:

- **GitHub Contributors**: Automatically listed on the repository
- **Release Notes**: Significant contributions mentioned in releases
- **Hall of Fame**: Outstanding contributors recognized in documentation

Thank you for contributing to flags-gen! Your efforts help make the Go ecosystem better for everyone.