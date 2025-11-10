package rules

import (
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// EmailRule generates email addresses.
type emailRule struct{}

// Email creates a rule that generates email addresses.
func Email() Rule {
	return &emailRule{}
}

func (r *emailRule) Generate(ctx Context) (interface{}, error) {
	domains := []string{"example.com", "test.com", "mail.com", "email.com"}
	prefixes := []string{"user", "test", "demo", "sample", "hello"}

	prefix := prefixes[ctx.Index()%len(prefixes)]
	domain := domains[ctx.Index()%len(domains)]

	return fmt.Sprintf("%s%d@%s", prefix, ctx.Index(), domain), nil
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func (r *emailRule) Validate(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", v)
	}
	if !emailRegex.MatchString(s) {
		return fmt.Errorf("invalid email format: %s", s)
	}
	return nil
}

// URLRule generates URLs.
type urlRule struct {
	scheme string
}

// URL creates a rule that generates URLs with the default scheme (https).
func URL() Rule {
	return &urlRule{scheme: "https"}
}

// URLWithScheme creates a rule that generates URLs with a specific scheme.
func URLWithScheme(scheme string) Rule {
	return &urlRule{scheme: scheme}
}

func (r *urlRule) Generate(ctx Context) (interface{}, error) {
	domains := []string{"example.com", "test.com", "demo.org", "sample.net"}
	paths := []string{"/", "/home", "/about", "/contact", "/products", "/services"}

	domain := domains[ctx.Index()%len(domains)]
	path := paths[ctx.Index()%len(paths)]

	return fmt.Sprintf("%s://%s%s", r.scheme, domain, path), nil
}

func (r *urlRule) Validate(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", v)
	}
	_, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	return nil
}

// UUIDRule generates UUIDs.
type uuidRule struct{}

// UUID creates a rule that generates UUID v4 strings.
func UUID() Rule {
	return &uuidRule{}
}

func (r *uuidRule) Generate(ctx Context) (interface{}, error) {
	// Use deterministic UUID based on seed and index for reproducibility
	seed := ctx.Seed() + int64(ctx.Index())
	rng := rand.New(rand.NewSource(seed))

	var uuidBytes [16]byte
	rng.Read(uuidBytes[:])

	// Set version (4) and variant bits
	uuidBytes[6] = (uuidBytes[6] & 0x0f) | 0x40
	uuidBytes[8] = (uuidBytes[8] & 0x3f) | 0x80

	return uuid.UUID(uuidBytes).String(), nil
}

func (r *uuidRule) Validate(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", v)
	}
	_, err := uuid.Parse(s)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}
	return nil
}

// RangeRule generates integers within a range.
type rangeRule struct {
	min int
	max int
}

// Range creates a rule that generates integers between min and max (inclusive).
func Range(min, max int) Rule {
	if min > max {
		min, max = max, min
	}
	return &rangeRule{min: min, max: max}
}

func (r *rangeRule) Generate(ctx Context) (interface{}, error) {
	if r.min == r.max {
		return r.min, nil
	}
	// Use index-based deterministic value
	rangeSize := r.max - r.min + 1
	return r.min + (ctx.Index() % rangeSize), nil
}

func (r *rangeRule) Validate(v interface{}) error {
	var n int
	switch val := v.(type) {
	case int:
		n = val
	case int64:
		n = int(val)
	case int32:
		n = int(val)
	default:
		return fmt.Errorf("expected integer type, got %T", v)
	}

	if n < r.min || n > r.max {
		return fmt.Errorf("value %d is outside range [%d, %d]", n, r.min, r.max)
	}
	return nil
}

// OneOfRule selects one value from a list of options.
type oneOfRule struct {
	options []interface{}
}

// OneOf creates a rule that randomly selects from the given options.
func OneOf(options ...interface{}) Rule {
	if len(options) == 0 {
		panic("OneOf requires at least one option")
	}
	return &oneOfRule{options: options}
}

func (r *oneOfRule) Generate(ctx Context) (interface{}, error) {
	idx := ctx.Index() % len(r.options)
	return r.options[idx], nil
}

func (r *oneOfRule) Validate(v interface{}) error {
	for _, opt := range r.options {
		if v == opt {
			return nil
		}
	}
	return fmt.Errorf("value %v is not one of the allowed options", v)
}

// SequenceRule generates sequential integers.
type sequenceRule struct {
	start int64
}

// Sequence creates a rule that generates sequential integers starting from start.
func Sequence(start int64) Rule {
	return &sequenceRule{start: start}
}

func (r *sequenceRule) Generate(ctx Context) (interface{}, error) {
	return r.start + int64(ctx.Index()), nil
}

func (r *sequenceRule) Validate(v interface{}) error {
	switch v.(type) {
	case int, int64, int32, int16, int8:
		return nil
	default:
		return fmt.Errorf("expected integer type, got %T", v)
	}
}

// AlphaNumericRule generates alphanumeric strings of a given length.
type alphaNumericRule struct {
	length int
}

// AlphaNumeric creates a rule that generates alphanumeric strings of the specified length.
func AlphaNumeric(length int) Rule {
	if length <= 0 {
		panic("AlphaNumeric length must be positive")
	}
	return &alphaNumericRule{length: length}
}

const alphaNumericChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (r *alphaNumericRule) Generate(ctx Context) (interface{}, error) {
	seed := ctx.Seed() + int64(ctx.Index())
	rng := rand.New(rand.NewSource(seed))

	var sb strings.Builder
	sb.Grow(r.length)

	for i := 0; i < r.length; i++ {
		sb.WriteByte(alphaNumericChars[rng.Intn(len(alphaNumericChars))])
	}

	return sb.String(), nil
}

func (r *alphaNumericRule) Validate(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", v)
	}
	if len(s) != r.length {
		return fmt.Errorf("expected length %d, got %d", r.length, len(s))
	}
	for _, c := range s {
		if !strings.ContainsRune(alphaNumericChars, c) {
			return fmt.Errorf("invalid character %q in alphanumeric string", c)
		}
	}
	return nil
}

// BoolRule generates boolean values with a specified probability of true.
type boolRule struct {
	trueRatio float64
}

// Bool creates a rule that generates boolean values.
// trueRatio is the probability (0.0 to 1.0) of generating true.
func Bool(trueRatio float64) Rule {
	if trueRatio < 0 || trueRatio > 1 {
		panic("Bool trueRatio must be between 0 and 1")
	}
	return &boolRule{trueRatio: trueRatio}
}

func (r *boolRule) Generate(ctx Context) (interface{}, error) {
	seed := ctx.Seed() + int64(ctx.Index())
	rng := rand.New(rand.NewSource(seed))
	return rng.Float64() < r.trueRatio, nil
}

func (r *boolRule) Validate(v interface{}) error {
	_, ok := v.(bool)
	if !ok {
		return fmt.Errorf("expected bool, got %T", v)
	}
	return nil
}

// DefaultRuleSet returns a RuleSet with all built-in rules registered.
func DefaultRuleSet() *RuleSet {
	rs := NewRuleSet()
	rs.Add("email", Email())
	rs.Add("url", URL())
	rs.Add("uuid", UUID())
	rs.Add("alphanumeric", AlphaNumeric(10))
	rs.Add("bool", Bool(0.5))
	return rs
}
