package rules

import (
	"errors"
	"testing"
)

type testRule struct {
	value interface{}
	err   error
}

func (r *testRule) Generate(ctx Context) (interface{}, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.value, nil
}

func (r *testRule) Validate(v interface{}) error {
	if v == r.value {
		return nil
	}
	return errors.New("validation failed")
}

func TestNewRuleSet(t *testing.T) {
	rs := NewRuleSet()
	if rs == nil {
		t.Fatal("NewRuleSet returned nil")
	}
	if rs.rules == nil {
		t.Error("rules map should be initialized")
	}
}

func TestRuleSet_Add(t *testing.T) {
	rs := NewRuleSet()
	rule := &testRule{value: "test"}

	rs.Add("test", rule)

	if !rs.Has("test") {
		t.Error("rule should be added")
	}

	retrieved, ok := rs.Get("test")
	if !ok {
		t.Error("should be able to get added rule")
	}
	if retrieved != rule {
		t.Error("retrieved rule should be the same as added rule")
	}
}

func TestRuleSet_AddChaining(t *testing.T) {
	rs := NewRuleSet()
	rule1 := &testRule{value: "test1"}
	rule2 := &testRule{value: "test2"}

	rs.Add("test1", rule1).Add("test2", rule2)

	if !rs.Has("test1") || !rs.Has("test2") {
		t.Error("both rules should be added")
	}
}

func TestRuleSet_Get(t *testing.T) {
	rs := NewRuleSet()
	rule := &testRule{value: "test"}
	rs.Add("test", rule)

	retrieved, ok := rs.Get("test")
	if !ok {
		t.Error("should find added rule")
	}
	if retrieved != rule {
		t.Error("retrieved rule should be the same as added rule")
	}

	_, ok = rs.Get("nonexistent")
	if ok {
		t.Error("should not find nonexistent rule")
	}
}

func TestRuleSet_Has(t *testing.T) {
	rs := NewRuleSet()
	rule := &testRule{value: "test"}
	rs.Add("test", rule)

	if !rs.Has("test") {
		t.Error("should have added rule")
	}
	if rs.Has("nonexistent") {
		t.Error("should not have nonexistent rule")
	}
}

func TestRuleSet_Remove(t *testing.T) {
	rs := NewRuleSet()
	rule := &testRule{value: "test"}
	rs.Add("test", rule)

	if !rs.Remove("test") {
		t.Error("should successfully remove existing rule")
	}
	if rs.Has("test") {
		t.Error("rule should be removed")
	}

	if rs.Remove("nonexistent") {
		t.Error("should return false for nonexistent rule")
	}
}

func TestRuleSet_Extend(t *testing.T) {
	rs1 := NewRuleSet()
	rs1.Add("rule1", &testRule{value: "test1"})
	rs1.Add("rule2", &testRule{value: "test2"})

	rs2 := NewRuleSet()
	rs2.Add("rule3", &testRule{value: "test3"})

	rs2.Extend(rs1)

	if !rs2.Has("rule1") || !rs2.Has("rule2") || !rs2.Has("rule3") {
		t.Error("all rules should be present after Extend")
	}
}

func TestRuleSet_ExtendOverwrite(t *testing.T) {
	rs1 := NewRuleSet()
	rule1 := &testRule{value: "original"}
	rs1.Add("test", rule1)

	rs2 := NewRuleSet()
	rule2 := &testRule{value: "new"}
	rs2.Add("test", rule2)

	rs2.Extend(rs1)

	retrieved, _ := rs2.Get("test")
	if retrieved != rule1 {
		t.Error("extended rule should overwrite existing rule")
	}
}

func TestRuleSet_ExtendNil(t *testing.T) {
	rs := NewRuleSet()
	rs.Add("test", &testRule{value: "test"})

	rs.Extend(nil)

	if !rs.Has("test") {
		t.Error("Extend(nil) should not affect existing rules")
	}
}

func TestRuleSet_Names(t *testing.T) {
	rs := NewRuleSet()
	rs.Add("rule1", &testRule{value: "test1"})
	rs.Add("rule2", &testRule{value: "test2"})
	rs.Add("rule3", &testRule{value: "test3"})

	names := rs.Names()
	if len(names) != 3 {
		t.Errorf("expected 3 names, got %d", len(names))
	}

	nameMap := make(map[string]bool)
	for _, name := range names {
		nameMap[name] = true
	}

	for _, expected := range []string{"rule1", "rule2", "rule3"} {
		if !nameMap[expected] {
			t.Errorf("expected to find %q in names", expected)
		}
	}
}

func TestRuleSet_Clone(t *testing.T) {
	rs1 := NewRuleSet()
	rule1 := &testRule{value: "test1"}
	rs1.Add("rule1", rule1)

	rs2 := rs1.Clone()

	// Check that clone has the same rules
	if !rs2.Has("rule1") {
		t.Error("clone should have the same rules")
	}

	// Check that modifications to clone don't affect original
	rs2.Add("rule2", &testRule{value: "test2"})
	if rs1.Has("rule2") {
		t.Error("modifications to clone should not affect original")
	}

	// Check that modifications to original don't affect clone
	rs1.Add("rule3", &testRule{value: "test3"})
	if rs2.Has("rule3") {
		t.Error("modifications to original should not affect clone")
	}
}

func TestRuleSet_ConcurrentAccess(t *testing.T) {
	rs := NewRuleSet()

	// Add initial rule
	rs.Add("rule1", &testRule{value: "test1"})

	done := make(chan bool)

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				rs.Get("rule1")
				rs.Has("rule1")
			}
			done <- true
		}()
	}

	// Concurrent writes
	for i := 0; i < 5; i++ {
		go func(id int) {
			for j := 0; j < 50; j++ {
				rs.Add("concurrent", &testRule{value: id})
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 15; i++ {
		<-done
	}
}

func TestValidationError(t *testing.T) {
	innerErr := errors.New("inner error")
	err := &ValidationError{
		RuleName: "test",
		Value:    "value",
		Err:      innerErr,
	}

	if err.Error() == "" {
		t.Error("Error() should return non-empty string")
	}

	if err.Unwrap() != innerErr {
		t.Error("Unwrap() should return inner error")
	}
}

func TestGenerationError(t *testing.T) {
	innerErr := errors.New("inner error")
	err := &GenerationError{
		RuleName: "test",
		Err:      innerErr,
	}

	if err.Error() == "" {
		t.Error("Error() should return non-empty string")
	}

	if err.Unwrap() != innerErr {
		t.Error("Unwrap() should return inner error")
	}
}
