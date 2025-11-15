package rules

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

// mockContext is a simple Context implementation for testing
type mockContext struct {
	locale string
	seed   int64
	index  int
}

func (m *mockContext) Locale() string                           { return m.locale }
func (m *mockContext) Seed() int64                              { return m.seed }
func (m *mockContext) Index() int                               { return m.index }
func (m *mockContext) GetField(name string) (interface{}, bool) { return nil, false }
func (m *mockContext) GetStruct() interface{}                   { return nil }
func (m *mockContext) FieldName() string                        { return "" }

func newMockContext(index int) Context {
	return &mockContext{
		locale: "en_US",
		seed:   12345,
		index:  index,
	}
}

func TestEmailRule(t *testing.T) {
	rule := Email()
	ctx := newMockContext(0)

	val, err := rule.Generate(ctx)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	email, ok := val.(string)
	if !ok {
		t.Fatalf("expected string, got %T", val)
	}

	if !strings.Contains(email, "@") {
		t.Errorf("expected email to contain @, got %s", email)
	}

	// Test validation
	if err := rule.Validate(email); err != nil {
		t.Errorf("Validate failed: %v", err)
	}

	// Test invalid email
	if err := rule.Validate("invalid-email"); err == nil {
		t.Error("expected validation error for invalid email")
	}
}

func TestURLRule(t *testing.T) {
	rule := URL()
	ctx := newMockContext(0)

	val, err := rule.Generate(ctx)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	url, ok := val.(string)
	if !ok {
		t.Fatalf("expected string, got %T", val)
	}

	if !strings.HasPrefix(url, "https://") {
		t.Errorf("expected URL to start with https://, got %s", url)
	}

	// Test validation
	if err := rule.Validate(url); err != nil {
		t.Errorf("Validate failed: %v", err)
	}
}

func TestURLWithScheme(t *testing.T) {
	rule := URLWithScheme("http")
	ctx := newMockContext(0)

	val, err := rule.Generate(ctx)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	url, ok := val.(string)
	if !ok {
		t.Fatalf("expected string, got %T", val)
	}

	if !strings.HasPrefix(url, "http://") {
		t.Errorf("expected URL to start with http://, got %s", url)
	}
}

func TestUUIDRule(t *testing.T) {
	rule := UUID()
	ctx := newMockContext(0)

	val, err := rule.Generate(ctx)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	uuidStr, ok := val.(string)
	if !ok {
		t.Fatalf("expected string, got %T", val)
	}

	// Test validation
	if err := rule.Validate(uuidStr); err != nil {
		t.Errorf("Validate failed: %v", err)
	}

	// Verify it's a valid UUID
	_, err = uuid.Parse(uuidStr)
	if err != nil {
		t.Errorf("invalid UUID generated: %v", err)
	}
}

func TestUUIDRule_Deterministic(t *testing.T) {
	rule := UUID()
	ctx1 := newMockContext(0)
	ctx2 := newMockContext(0)

	val1, _ := rule.Generate(ctx1)
	val2, _ := rule.Generate(ctx2)

	if val1 != val2 {
		t.Errorf("expected same UUID for same context, got %v and %v", val1, val2)
	}
}

func TestRangeRule(t *testing.T) {
	rule := Range(10, 20)

	for i := 0; i < 100; i++ {
		ctx := newMockContext(i)
		val, err := rule.Generate(ctx)
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		n, ok := val.(int)
		if !ok {
			t.Fatalf("expected int, got %T", val)
		}

		if n < 10 || n > 20 {
			t.Errorf("value %d is outside range [10, 20]", n)
		}

		// Test validation
		if err := rule.Validate(n); err != nil {
			t.Errorf("Validate failed for %d: %v", n, err)
		}
	}

	// Test validation for out of range value
	if err := rule.Validate(5); err == nil {
		t.Error("expected validation error for value outside range")
	}
}

func TestRangeRule_SwappedMinMax(t *testing.T) {
	rule := Range(20, 10) // min > max
	ctx := newMockContext(0)

	val, err := rule.Generate(ctx)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	n := val.(int)
	if n < 10 || n > 20 {
		t.Errorf("value %d is outside range [10, 20]", n)
	}
}

func TestOneOfRule(t *testing.T) {
	options := []interface{}{"red", "green", "blue"}
	rule := OneOf(options...)

	seen := make(map[string]bool)
	for i := 0; i < 10; i++ {
		ctx := newMockContext(i)
		val, err := rule.Generate(ctx)
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		str, ok := val.(string)
		if !ok {
			t.Fatalf("expected string, got %T", val)
		}

		found := false
		for _, opt := range options {
			if opt == str {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("value %s is not in options", str)
		}

		seen[str] = true

		// Test validation
		if err := rule.Validate(val); err != nil {
			t.Errorf("Validate failed: %v", err)
		}
	}

	// Test validation for invalid value
	if err := rule.Validate("yellow"); err == nil {
		t.Error("expected validation error for value not in options")
	}
}

func TestOneOfRule_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for OneOf with no options")
		}
	}()
	OneOf()
}

func TestSequenceRule(t *testing.T) {
	rule := Sequence(100)

	for i := 0; i < 10; i++ {
		ctx := newMockContext(i)
		val, err := rule.Generate(ctx)
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		expected := int64(100 + i)
		if val != expected {
			t.Errorf("expected %d, got %v", expected, val)
		}

		// Test validation
		if err := rule.Validate(val); err != nil {
			t.Errorf("Validate failed: %v", err)
		}
	}
}

func TestAlphaNumericRule(t *testing.T) {
	rule := AlphaNumeric(10)
	ctx := newMockContext(0)

	val, err := rule.Generate(ctx)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	str, ok := val.(string)
	if !ok {
		t.Fatalf("expected string, got %T", val)
	}

	if len(str) != 10 {
		t.Errorf("expected length 10, got %d", len(str))
	}

	// Verify all characters are alphanumeric
	for _, c := range str {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			t.Errorf("non-alphanumeric character %c in %s", c, str)
		}
	}

	// Test validation
	if err := rule.Validate(str); err != nil {
		t.Errorf("Validate failed: %v", err)
	}
}

func TestAlphaNumericRule_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for AlphaNumeric with non-positive length")
		}
	}()
	AlphaNumeric(0)
}

func TestBoolRule(t *testing.T) {
	rule := Bool(0.5)

	trueCount := 0
	falseCount := 0

	for i := 0; i < 100; i++ {
		ctx := newMockContext(i)
		val, err := rule.Generate(ctx)
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		b, ok := val.(bool)
		if !ok {
			t.Fatalf("expected bool, got %T", val)
		}

		if b {
			trueCount++
		} else {
			falseCount++
		}

		// Test validation
		if err := rule.Validate(val); err != nil {
			t.Errorf("Validate failed: %v", err)
		}
	}

	// With 0.5 ratio, we should get some of both
	if trueCount == 0 || falseCount == 0 {
		t.Errorf("expected mix of true and false, got %d true and %d false", trueCount, falseCount)
	}
}

func TestBoolRule_Panic(t *testing.T) {
	tests := []float64{-0.1, 1.1}
	for _, ratio := range tests {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic for Bool with ratio %f", ratio)
				}
			}()
			Bool(ratio)
		}()
	}
}

func TestDefaultRuleSet(t *testing.T) {
	rs := DefaultRuleSet()

	expectedRules := []string{"email", "url", "uuid", "alphanumeric", "bool"}
	for _, name := range expectedRules {
		if !rs.Has(name) {
			t.Errorf("expected rule %q to be in default RuleSet", name)
		}
	}
}

func BenchmarkEmailRule(b *testing.B) {
	rule := Email()
	ctx := newMockContext(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.Generate(ctx)
	}
}

func BenchmarkUUIDRule(b *testing.B) {
	rule := UUID()
	ctx := newMockContext(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.Generate(ctx)
	}
}

func BenchmarkRangeRule(b *testing.B) {
	rule := Range(1, 100)
	ctx := newMockContext(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.Generate(ctx)
	}
}
