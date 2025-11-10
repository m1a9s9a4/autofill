package rules

import (
	"fmt"
	"sync"
)

// Context provides information about the current generation context.
// This is re-exported from the parent package to avoid circular dependencies.
type Context interface {
	Locale() string
	Seed() int64
	Index() int
	GetField(name string) (interface{}, bool)
	GetStruct() interface{}
	FieldName() string
}

// Rule defines the interface for value generation rules.
// Rules can generate values based on context and validate generated values.
type Rule interface {
	// Generate creates a new value based on the given context.
	Generate(ctx Context) (interface{}, error)

	// Validate checks if the given value is valid for this rule.
	// Returns nil if valid, or an error describing the validation failure.
	Validate(v interface{}) error
}

// LocalizedRule is a Rule that supports multiple locales.
type LocalizedRule interface {
	Rule

	// SupportedLocales returns a list of locale codes this rule supports.
	// Examples: "ja_JP", "en_US", "ko_KR"
	SupportedLocales() []string
}

// RuleSet manages a collection of named rules.
// It is safe for concurrent use.
type RuleSet struct {
	mu    sync.RWMutex
	rules map[string]Rule
}

// NewRuleSet creates a new empty RuleSet.
func NewRuleSet() *RuleSet {
	return &RuleSet{
		rules: make(map[string]Rule),
	}
}

// Add registers a rule with the given name.
// If a rule with the same name already exists, it will be replaced.
// Returns the RuleSet for method chaining.
func (rs *RuleSet) Add(name string, rule Rule) *RuleSet {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.rules[name] = rule
	return rs
}

// Get retrieves a rule by name.
// Returns the rule and true if found, or nil and false if not found.
func (rs *RuleSet) Get(name string) (Rule, bool) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	rule, ok := rs.rules[name]
	return rule, ok
}

// Has checks if a rule with the given name exists.
func (rs *RuleSet) Has(name string) bool {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	_, ok := rs.rules[name]
	return ok
}

// Remove removes a rule by name.
// Returns true if the rule was found and removed, false otherwise.
func (rs *RuleSet) Remove(name string) bool {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if _, ok := rs.rules[name]; ok {
		delete(rs.rules, name)
		return true
	}
	return false
}

// Extend adds all rules from another RuleSet to this one.
// Existing rules with the same name will be replaced.
// Returns the RuleSet for method chaining.
func (rs *RuleSet) Extend(other *RuleSet) *RuleSet {
	if other == nil {
		return rs
	}

	other.mu.RLock()
	defer other.mu.RUnlock()

	rs.mu.Lock()
	defer rs.mu.Unlock()

	for name, rule := range other.rules {
		rs.rules[name] = rule
	}
	return rs
}

// Names returns a list of all registered rule names.
func (rs *RuleSet) Names() []string {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	names := make([]string, 0, len(rs.rules))
	for name := range rs.rules {
		names = append(names, name)
	}
	return names
}

// Clone creates a shallow copy of the RuleSet.
// The rules themselves are not cloned, only the map structure.
func (rs *RuleSet) Clone() *RuleSet {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	newRS := NewRuleSet()
	for name, rule := range rs.rules {
		newRS.rules[name] = rule
	}
	return newRS
}

// ValidationError represents a validation error from a rule.
type ValidationError struct {
	RuleName string
	Value    interface{}
	Err      error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for rule %q with value %v: %v", e.RuleName, e.Value, e.Err)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// GenerationError represents an error during value generation.
type GenerationError struct {
	RuleName string
	Err      error
}

func (e *GenerationError) Error() string {
	return fmt.Sprintf("generation failed for rule %q: %v", e.RuleName, e.Err)
}

func (e *GenerationError) Unwrap() error {
	return e.Err
}
