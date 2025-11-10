# autofill

[![Go Reference](https://pkg.go.dev/badge/github.com/m1a9s9a4/autofill.svg)](https://pkg.go.dev/github.com/m1a9s9a4/autofill)
[![Go Report Card](https://goreportcard.com/badge/github.com/m1a9s9a4/autofill)](https://goreportcard.com/report/github.com/m1a9s9a4/autofill)
[![Coverage](https://codecov.io/gh/m1a9s9a4/autofill/branch/main/graph/badge.svg)](https://codecov.io/gh/m1a9s9a4/autofill)

Automatically fill Go structs with realistic test data.

## Features

- 🚀 **Zero dependencies** - Uses only the Go standard library (except for UUID generation)
- 🎯 **Type-safe** - Compile-time type checking with struct tags
- 🔧 **Extensible** - Easy to add custom rules
- 📝 **Simple API** - Clean and intuitive interface
- 🧪 **Well-tested** - 80%+ test coverage
- ⚡ **Fast** - Optimized for performance
- 🎲 **Deterministic** - Same seed produces same results

## Installation

```bash
go get github.com/m1a9s9a4/autofill
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/m1a9s9a4/autofill"
)

type User struct {
    ID    int64  `autofill:"seq"`
    Name  string
    Email string `autofill:"email"`
    Age   int    `autofill:"min=18,max=65"`
}

func main() {
    var user User
    autofill.Fill(&user)
    fmt.Printf("%+v\n", user)
    // Output: {ID:0 Name:hello Email:user0@example.com Age:18}
}
```

## Usage

### Basic Fill

Fill a single struct with generated data:

```go
var user User
err := autofill.Fill(&user)
if err != nil {
    panic(err)
}
```

### Fill Multiple Structs

Fill a slice of structs:

```go
users := make([]User, 10)
err := autofill.FillSlice(&users)
if err != nil {
    panic(err)
}
```

### Configuration

Configure autofill with method chaining:

```go
af := autofill.New().
    WithSeed(12345).        // Deterministic generation
    WithLocale("ja_JP")     // Set locale (for future locale-aware rules)

af.Fill(&user)
```

### Override Values

Override specific field values:

```go
autofill.Fill(&user, autofill.Override{
    "Name":  "John Doe",
    "Email": "john@example.com",
    "Age":   30,
})
```

### Sequence Values

Generate sequential values for slices:

```go
users := make([]User, 5)
autofill.FillSlice(&users, autofill.Override{
    "Email": autofill.Seq("user%d@example.com"),  // user0@, user1@, ...
    "Age":   autofill.SeqInt(20),                  // 20, 21, 22, ...
    "ID":    autofill.SeqInt64(1000),              // 1000, 1001, 1002, ...
})
```

### Struct Tags

Use struct tags to control value generation:

```go
type User struct {
    ID        int64     `autofill:"seq"`              // Sequential: 0, 1, 2, ...
    Email     string    `autofill:"email"`            // Email: user0@example.com
    URL       string    `autofill:"url"`              // URL: https://example.com/
    UUID      string    `autofill:"uuid"`             // UUID: v4 format
    Age       int       `autofill:"min=18,max=65"`    // Range: 18-65
    Status    string    `autofill:"oneof=active|inactive"` // One of: active or inactive
    CreatedAt time.Time `autofill:"now"`              // Current time
    Internal  string    `autofill:"-"`                // Skip this field
}
```

### Available Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `seq` | Sequential integers starting from 0 | `autofill:"seq"` |
| `email` | Email addresses | `autofill:"email"` |
| `url` | URLs | `autofill:"url"` |
| `uuid` | UUID v4 strings | `autofill:"uuid"` |
| `now` | Current time | `autofill:"now"` |
| `min=N,max=M` | Integer range [N, M] | `autofill:"min=18,max=65"` |
| `oneof=a\|b\|c` | Choose from options | `autofill:"oneof=active\|inactive"` |
| `rule=name` | Use custom rule | `autofill:"rule=myRule"` |
| `-` | Skip field | `autofill:"-"` |

### Custom Rules

Create and register custom rules:

```go
package main

import (
    "github.com/m1a9s9a4/autofill"
    "github.com/m1a9s9a4/autofill/rules"
)

// Custom rule implementation
type StatusRule struct{}

func (r *StatusRule) Generate(ctx rules.Context) (interface{}, error) {
    statuses := []string{"active", "inactive", "pending"}
    return statuses[ctx.Index()%len(statuses)], nil
}

func (r *StatusRule) Validate(v interface{}) error {
    return nil
}

func main() {
    // Create RuleSet and add custom rule
    ruleSet := rules.DefaultRuleSet()
    ruleSet.Add("status", &StatusRule{})

    type Task struct {
        Status string `autofill:"rule=status"`
    }

    // Use custom rule
    af := autofill.New().WithRules(ruleSet)
    var task Task
    af.Fill(&task)
}
```

### Testing

Use autofill in your tests:

```go
func TestUserCreation(t *testing.T) {
    var user User
    if err := autofill.Fill(&user); err != nil {
        t.Fatalf("Fill failed: %v", err)
    }

    if user.Email == "" {
        t.Error("Email should not be empty")
    }
}

func TestBulkUsers(t *testing.T) {
    users := make([]User, 100)
    if err := autofill.FillSlice(&users); err != nil {
        t.Fatalf("FillSlice failed: %v", err)
    }

    // All users should have different IDs
    for i, user := range users {
        if user.ID != int64(i) {
            t.Errorf("user %d: expected ID %d, got %d", i, i, user.ID)
        }
    }
}
```

## Built-in Rules

The package includes several built-in rules accessible via tags or the rules API:

### String Rules
- **Email**: Generates email addresses (e.g., `user0@example.com`)
- **URL**: Generates URLs (e.g., `https://example.com/`)
- **UUID**: Generates UUID v4 strings
- **AlphaNumeric**: Generates alphanumeric strings of specified length

### Numeric Rules
- **Range**: Generates integers within a range
- **Sequence**: Generates sequential integers

### Selection Rules
- **OneOf**: Selects from a list of options

### Other Rules
- **Bool**: Generates boolean values with configurable probability

## Examples

See the [examples](examples/) directory for complete examples:

- [Basic Usage](examples/basic/main.go) - Basic filling and overrides
- [Custom Rules](examples/custom_rules/main.go) - Creating custom rules
- [Testing](examples/testing/user_test.go) - Using autofill in tests

## API Documentation

Full API documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/m1a9s9a4/autofill).

### Main Functions

```go
// Create new Autofill instance
func New() *Autofill

// Configure Autofill
func (a *Autofill) WithLocale(locale string) *Autofill
func (a *Autofill) WithSeed(seed int64) *Autofill
func (a *Autofill) WithRules(rules *RuleSet) *Autofill

// Fill structs
func (a *Autofill) Fill(v interface{}, overrides ...Override) error
func (a *Autofill) FillSlice(v interface{}, overrides ...Override) error

// Convenience functions
func Fill(v interface{}, overrides ...Override) error
func FillSlice(v interface{}, overrides ...Override) error
```

### Override Functions

```go
// Create sequence functions for overrides
func Seq(format string) SequenceFunc
func SeqInt(start int) SequenceFunc
func SeqInt64(start int64) SequenceFunc
func Random(min, max int) SequenceFunc
```

## Supported Types

autofill supports the following Go types:

- **Basic types**: `string`, `bool`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`
- **Time**: `time.Time`
- **Pointers**: Pointers to any supported type
- **Structs**: Nested struct types
- **Slices**: Slices of any supported type

## Performance

autofill is optimized for performance:

```
BenchmarkFill-8              500000    2341 ns/op     896 B/op    18 allocs/op
BenchmarkFillSlice-8           5000  234156 ns/op   89600 B/op  1800 allocs/op
```

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- Inspired by libraries like Faker and Go's testing/quick package
- Built with ❤️ for the Go community
