# flags-gen

[![Go Version](https://img.shields.io/github/go-mod/go-version/yuvalwz/flags-gen)](https://golang.org/releases)
[![License](https://img.shields.io/github/license/yuvalwz/flags-gen)](LICENSE)
[![Release](https://img.shields.io/github/v/release/yuvalwz/flags-gen)](https://github.com/yuvalwz/flags-gen/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yuvalwz/flags-gen)](https://goreportcard.com/report/github.com/yuvalwz/flags-gen)

A powerful Go code generation tool that automatically creates [pflags](https://github.com/spf13/pflag) `AddFlags` methods from Go structs with simple annotations. Perfect for CLI applications using [Cobra](https://github.com/spf13/cobra) or any application that needs command-line flag parsing.

## Features

- **Simple Annotation**: Just add `+flags-gen` comment above your struct
- **Type Safety**: Generates strongly-typed flag methods
- **Smart Naming**: Converts camelCase field names to kebab-case flags
- **Default Values**: Uses struct tags for default values
- **Rich Types**: Supports strings, integers, booleans, slices, durations, and more
- **Documentation**: Extracts flag descriptions from Go comments
- **Zero Dependencies**: Generated code only depends on `pflag`

## Supported Types

| Go Type | Flag Method | Example Flag |
|---------|-------------|--------------|
| `string` | `StringVar` | `--name value` |
| `int`, `int32`, `int64` | `IntVar`, `Int32Var`, `Int64Var` | `--count 42` |
| `uint`, `uint32`, `uint64` | `UintVar`, `Uint32Var`, `Uint64Var` | `--size 1024` |
| `bool` | `BoolVar` | `--enabled` |
| `float32`, `float64` | `Float32Var`, `Float64Var` | `--rate 0.5` |
| `[]string` | `StringSliceVar` | `--tags foo,bar` |
| `[]int` | `IntSliceVar` | `--ports 80,443` |
| `time.Duration` | `DurationVar` | `--timeout 30s` |

## Installation

### Using `go install`

```bash
go install github.com/yuvalwz/flags-gen/cmd/flags-gen@latest
```

### Using Homebrew

```bash
# Coming soon
brew install flags-gen
```

### Binary Releases

Download pre-built binaries from the [releases page](https://github.com/yuvalwz/flags-gen/releases).

### Building from Source

```bash
git clone https://github.com/yuvalwz/flags-gen.git
cd flags-gen
make build
```

## Quick Start

### 1. Annotate Your Struct

Add the `+flags-gen` annotation above your struct:

```go
package config

import "time"

// +flags-gen
// ServerConfig holds all server configuration options
type ServerConfig struct {
    // Host is the server host address
    Host string `json:"host" default:"localhost"`
    
    // Port is the server port number
    Port int `json:"port" default:"8080"`
    
    // EnableTLS enables HTTPS support
    EnableTLS bool `json:"enableTLS"`
    
    // AllowedOrigins contains allowed CORS origins
    AllowedOrigins []string `json:"allowedOrigins" default:"http://localhost:3000"`
    
    // RequestTimeout is the timeout for incoming requests
    RequestTimeout time.Duration `json:"requestTimeout" default:"30s"`
}
```

### 2. Generate Flags Code

```bash
flags-gen -i config.go -o config_flags.go
```

### 3. Use Generated Code

The tool generates a method you can use with any `pflag.FlagSet`:

```go
package main

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/spf13/pflag"
    "your-project/config"
)

func main() {
    cfg := &config.ServerConfig{}
    
    rootCmd := &cobra.Command{
        Use: "myapp",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("Server will run on %s:%d\n", cfg.Host, cfg.Port)
            if cfg.EnableTLS {
                fmt.Println("TLS is enabled")
            }
        },
    }
    
    // Add all flags from the struct
    cfg.AddFlags(rootCmd.Flags())
    
    rootCmd.Execute()
}
```

### 4. Use Your CLI

```bash
./myapp --host 0.0.0.0 --port 9000 --enable-tls --allowed-origins "https://example.com,https://app.com" --request-timeout 45s
```

## Usage

### Command Line Options

```bash
flags-gen -i <input-file> [-o <output-file>]
```

**Options:**
- `-i, --input`: Input Go file containing structs with `+flags-gen` annotations (required)
- `-o, --output`: Output file for generated flags code (optional, defaults to `<input>_flags.go`)
- `--version`: Show version information

**Examples:**
```bash
# Generate flags for types.go, output to types_flags.go
flags-gen -i types.go

# Generate flags with custom output file
flags-gen -i pkg/config/server.go -o pkg/config/flags.go

# Generate flags for multiple structs in one file
flags-gen -i internal/config/config.go -o internal/config/generated_flags.go
```

### Struct Tag Options

Control flag generation with struct tags:

```go
type Config struct {
    // Use json tag for custom flag name
    DatabaseURL string `json:"database-url" default:"sqlite://app.db"`
    
    // Set default values
    MaxConnections int `json:"maxConnections" default:"10"`
    
    // Boolean flags don't need defaults (default to false)
    DebugMode bool `json:"debugMode"`
    
    // String slices support comma-separated defaults
    Features []string `json:"features" default:"auth,logging,metrics"`
    
    // Duration fields support time string defaults
    CacheTimeout time.Duration `json:"cacheTimeout" default:"5m"`
}
```

### Comment-Based Documentation

The tool extracts flag descriptions from Go comments:

```go
type APIConfig struct {
    // APIKey is the secret key for API authentication.
    // This should be kept secure and never logged.
    APIKey string `json:"apiKey"`
    
    // RateLimit sets the maximum requests per second.
    // Set to 0 to disable rate limiting.
    RateLimit int `json:"rateLimit" default:"100"`
}
```

## Integration Examples

### With Cobra CLI

```go
package main

import (
    "github.com/spf13/cobra"
    "your-project/internal/config"
)

func main() {
    cfg := &config.Config{}
    
    rootCmd := &cobra.Command{
        Use:   "myapp",
        Short: "My awesome CLI application",
        RunE: func(cmd *cobra.Command, args []string) error {
            return runServer(cfg)
        },
    }
    
    // Add generated flags
    cfg.AddFlags(rootCmd.Flags())
    
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

### With Custom FlagSet

```go
package main

import (
    "flag"
    "github.com/spf13/pflag"
    "your-project/config"
)

func main() {
    cfg := &config.ServerConfig{}
    
    // Create custom flag set
    fs := pflag.NewFlagSet("myapp", pflag.ExitOnError)
    
    // Add generated flags
    cfg.AddFlags(fs)
    
    // Parse flags
    fs.Parse(os.Args[1:])
    
    // Use configuration
    startServer(cfg)
}
```

### With Environment Variable Support

Combine with [viper](https://github.com/spf13/viper) for environment variable support:

```go
package main

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "your-project/config"
)

func main() {
    cfg := &config.Config{}
    
    rootCmd := &cobra.Command{
        Use: "myapp",
        PreRunE: func(cmd *cobra.Command, args []string) error {
            // Bind environment variables
            viper.SetEnvPrefix("MYAPP")
            viper.AutomaticEnv()
            viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
            
            return viper.BindPFlags(cmd.Flags())
        },
        RunE: func(cmd *cobra.Command, args []string) error {
            return viper.Unmarshal(cfg)
        },
    }
    
    cfg.AddFlags(rootCmd.Flags())
    rootCmd.Execute()
}
```

## Advanced Features

### Multiple Structs

Generate flags for multiple structs in a single file:

```go
// +flags-gen
type ServerConfig struct {
    Host string `json:"host" default:"localhost"`
    Port int `json:"port" default:"8080"`
}

// +flags-gen  
type DatabaseConfig struct {
    URL string `json:"url" default:"sqlite://app.db"`
    MaxConns int `json:"maxConns" default:"10"`
}
```

Both structs will get their own `AddFlags` methods.

### Embedded Structs

The tool skips embedded struct fields to avoid conflicts:

```go
// +flags-gen
type Config struct {
    metav1.TypeMeta `json:",inline"` // Skipped
    
    // Regular fields are processed
    Name string `json:"name"`
    Port int    `json:"port" default:"8080"`
}
```

### Complex Default Values

Handle complex default values in various formats:

```go
type Config struct {
    // String slices
    Tags []string `json:"tags" default:"tag1,tag2,tag3"`
    
    // Duration with units
    Timeout time.Duration `json:"timeout" default:"30s"`
    
    // Numeric values
    BufferSize int `json:"bufferSize" default:"1024"`
    
    // Boolean values
    Enabled bool `json:"enabled" default:"true"`
}
```

## Development

### Prerequisites

- Go 1.20 or later
- Make (optional, for convenience commands)

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install to $GOPATH/bin
make install
```

### Testing

```bash
# Run tests
make test

# Run tests with race detection
make test-race

# Generate coverage report
make coverage
```

### Code Quality

```bash
# Format code
make fmt

# Run linting
make lint

# Run all checks
make check
```

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on:

- Setting up the development environment
- Code style and standards
- Testing requirements
- Submitting pull requests

## Examples

For more detailed examples and use cases, see [docs/examples.md](docs/examples.md).

## Troubleshooting

### Common Issues

**Issue**: Generated code has compilation errors
**Solution**: Ensure your struct uses supported types and has proper Go syntax

**Issue**: Flags not appearing in CLI
**Solution**: Make sure you're calling the `AddFlags` method on your flag set

**Issue**: Default values not working
**Solution**: Check that your struct tags use the correct `default:"value"` format

**Issue**: Flag names not as expected
**Solution**: Use `json` tags to override automatic kebab-case conversion

### Getting Help

- Check the [examples](docs/examples.md)
- Search [existing issues](https://github.com/yuvalwz/flags-gen/issues)
- Create a [new issue](https://github.com/yuvalwz/flags-gen/issues/new)

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [spf13/pflag](https://github.com/spf13/pflag) - POSIX/GNU-style command-line flags
- [spf13/cobra](https://github.com/spf13/cobra) - Modern CLI framework for Go
- Go team for the excellent `go/ast` and `go/parser` packages

---

Made with ❤️ for the Go community