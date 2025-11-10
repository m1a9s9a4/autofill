# Contributing to autofill

Thank you for your interest in contributing to autofill! This document provides guidelines and instructions for contributing.

## Code of Conduct

Be respectful, constructive, and professional in all interactions.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:

- A clear, descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Go version and OS
- Minimal code example demonstrating the issue

### Suggesting Features

Feature suggestions are welcome! Please create an issue with:

- A clear description of the feature
- Use cases for the feature
- Example API usage (if applicable)

### Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Write or update tests
5. Ensure all tests pass and coverage stays above 80%
6. Run linters and format code
7. Commit your changes with clear messages
8. Push to your fork
9. Create a Pull Request

## Development Setup

### Prerequisites

- Go 1.21 or later
- golangci-lint (for linting)

### Clone and Setup

```bash
git clone https://github.com/m1a9s9a4/autofill.git
cd autofill
go mod download
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -run TestName ./...

# Run benchmarks
go test -bench=. -benchmem ./...
```

### Code Coverage Requirements

- **Minimum coverage**: 80%
- Check coverage: `go tool cover -func=coverage.out | grep total`
- Coverage must not decrease with new contributions

### Linting

```bash
# Run linter
golangci-lint run

# Auto-fix issues where possible
golangci-lint run --fix
```

### Code Formatting

```bash
# Format code
go fmt ./...

# Organize imports
goimports -w .
```

## Adding New Rules

New rules are highly encouraged! Here's how to add one:

### 1. Create Rule File

Create `rules/yourrule.go`:

```go
package rules

import "fmt"

// YourRule generates your custom values
type yourRule struct {
    // configuration fields
}

// YourRule creates a new YourRule
func YourRule() Rule {
    return &yourRule{}
}

func (r *yourRule) Generate(ctx Context) (interface{}, error) {
    // Your generation logic here
    return "generated value", nil
}

func (r *yourRule) Validate(v interface{}) error {
    // Validation logic
    if v == nil {
        return fmt.Errorf("value cannot be nil")
    }
    return nil
}
```

### 2. Write Tests

Create `rules/yourrule_test.go`:

```go
package rules

import "testing"

func TestYourRule(t *testing.T) {
    rule := YourRule()
    ctx := newMockContext(0)

    // Test generation
    val, err := rule.Generate(ctx)
    if err != nil {
        t.Fatalf("Generate failed: %v", err)
    }

    // Verify value
    str, ok := val.(string)
    if !ok {
        t.Fatalf("expected string, got %T", val)
    }
    if str == "" {
        t.Error("generated value should not be empty")
    }

    // Test validation
    if err := rule.Validate(val); err != nil {
        t.Errorf("Validate failed: %v", err)
    }
}

func TestYourRule_Deterministic(t *testing.T) {
    rule := YourRule()

    ctx1 := newMockContext(0)
    val1, _ := rule.Generate(ctx1)

    ctx2 := newMockContext(0)
    val2, _ := rule.Generate(ctx2)

    if val1 != val2 {
        t.Errorf("expected same value for same context")
    }
}

func BenchmarkYourRule(b *testing.B) {
    rule := YourRule()
    ctx := newMockContext(0)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        rule.Generate(ctx)
    }
}
```

### 3. Required Tests

Every rule MUST have:

1. **Basic generation test** - Verify it generates valid values
2. **Type test** - Verify correct return type
3. **Validation test** - Test the Validate method
4. **Deterministic test** - Same context produces same results
5. **Edge cases** - Test boundary conditions
6. **Benchmark** - Performance test

### 4. Update DefaultRuleSet (Optional)

If the rule should be included by default, add it to `DefaultRuleSet()` in `rules/builtin.go`:

```go
func DefaultRuleSet() *RuleSet {
    rs := NewRuleSet()
    rs.Add("email", Email())
    rs.Add("url", URL())
    rs.Add("uuid", UUID())
    rs.Add("your-rule", YourRule())  // Add your rule here
    return rs
}
```

### 5. Document Your Rule

Add documentation to README.md:

```markdown
### YourRule

Generates custom values...

Usage:
\`\`\`go
type Example struct {
    Field string `autofill:"rule=yourrule"`
}
\`\`\`
```

## Localized Rules

If your rule should support multiple locales, implement the `LocalizedRule` interface:

```go
type yourLocalizedRule struct {
    data map[string]*YourData
}

func (r *yourLocalizedRule) SupportedLocales() []string {
    return []string{"ja_JP", "en_US", "ko_KR"}
}

func (r *yourLocalizedRule) Generate(ctx Context) (interface{}, error) {
    locale := ctx.Locale()
    data := r.data[locale]
    if data == nil {
        data = r.data["en_US"] // fallback
    }
    // Generate using locale-specific data
}
```

### Adding Locale Data

1. Create directory: `locale/[locale_code]/`
2. Add JSON file: `locale/[locale_code]/yourdata.json`
3. Load in your rule's init or constructor

Example:

```json
{
  "items": ["item1", "item2", "item3"],
  "categories": ["cat1", "cat2"]
}
```

## Code Style

### General Guidelines

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use meaningful variable and function names
- Keep functions small and focused
- Add comments for exported types and functions
- Use Go idioms and conventions

### Comments

```go
// Good: Concise and clear
// Email generates email addresses in the format user@domain.com
func Email() Rule {
    return &emailRule{}
}

// Bad: Too verbose or obvious
// Email is a function that returns a Rule interface implementation
// which can be used to generate email address strings
func Email() Rule {
    return &emailRule{}
}
```

### Error Handling

```go
// Good: Wrap errors with context
if err != nil {
    return nil, fmt.Errorf("failed to generate email: %w", err)
}

// Bad: Lose error context
if err != nil {
    return nil, err
}
```

## Commit Messages

Use conventional commit format:

- `feat: add new rule for phone numbers`
- `fix: correct email validation regex`
- `docs: update README with examples`
- `test: add tests for range rule`
- `refactor: simplify context creation`
- `perf: optimize UUID generation`
- `chore: update dependencies`

## PR Guidelines

### Before Submitting

- [ ] All tests pass (`go test ./...`)
- [ ] Coverage is above 80% (`go tool cover -func=coverage.out`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] Documentation is updated
- [ ] Examples are added/updated if needed

### PR Title

Use conventional commit format:

- `feat: add phone number rule`
- `fix: handle nil pointers correctly`
- `docs: improve custom rules documentation`

### PR Description

Include:

- **What**: What does this PR do?
- **Why**: Why is this change needed?
- **How**: How does it work?
- **Testing**: How was it tested?
- **Breaking Changes**: Any breaking changes? (mark as BREAKING CHANGE)

Example:

```markdown
## What
Adds a new rule for generating phone numbers with locale support.

## Why
Users need realistic phone number generation for testing.

## How
- Implemented PhoneRule with LocalizedRule interface
- Added phone number data for ja_JP, en_US, ko_KR
- Supports various formats based on locale

## Testing
- Added unit tests with 95% coverage
- Added benchmarks
- Tested with all supported locales
- Added example in examples/phone/

## Breaking Changes
None
```

## Questions?

- Open an issue for questions
- Check existing issues and PRs
- Read the [README](README.md) and examples

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
