// Package autofill provides automatic test data generation for Go structs.
//
// It supports:
//   - Automatic type-based value generation
//   - Custom rules via struct tags
//   - Locale-aware data generation
//   - Override values and sequences
//   - Deterministic generation with seeds
//
// Basic usage:
//
//	type User struct {
//	    Name  string `autofill:"rule=name"`
//	    Email string `autofill:"rule=email"`
//	    Age   int    `autofill:"min=18,max=65"`
//	}
//
//	var user User
//	autofill.Fill(&user)
//
// For more examples, see the examples/ directory in the repository.
package autofill

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/m1a9s9a4/autofill/rules"
)

// Autofill is the main struct for generating test data.
// Create an instance using New() and configure it with With* methods.
type Autofill struct {
	locale string
	seed   int64
	rules  *rules.RuleSet
	rand   *rand.Rand
}

// New creates a new Autofill instance with default settings.
// Default locale is "en_US" and seed is based on current time.
func New() *Autofill {
	seed := time.Now().UnixNano()
	return &Autofill{
		locale: "en_US",
		seed:   seed,
		rules:  rules.DefaultRuleSet(),
		rand:   rand.New(rand.NewSource(seed)),
	}
}

// WithLocale sets the locale for data generation.
// Common locales: "ja_JP" (Japanese), "en_US" (English), "ko_KR" (Korean).
func (a *Autofill) WithLocale(locale string) *Autofill {
	a.locale = locale
	return a
}

// WithSeed sets the random seed for deterministic generation.
// Using the same seed will produce the same results.
func (a *Autofill) WithSeed(seed int64) *Autofill {
	a.seed = seed
	a.rand = rand.New(rand.NewSource(seed))
	return a
}

// WithRules sets a custom RuleSet for value generation.
// This replaces the default RuleSet. Use Extend() to add to existing rules.
func (a *Autofill) WithRules(ruleSet *rules.RuleSet) *Autofill {
	a.rules = ruleSet
	return a
}

// Fill populates the fields of a struct with generated test data.
// The input must be a pointer to a struct.
// Optional overrides can be provided to set specific field values.
//
// Example:
//
//	var user User
//	err := autofill.New().Fill(&user, autofill.Override{
//	    "Name": "John Doe",
//	    "Age":  30,
//	})
func (a *Autofill) Fill(v interface{}, overrides ...Override) error {
	return a.FillWithIndex(v, 0, overrides...)
}

// FillWithIndex is like Fill but allows specifying a custom index.
// This is useful for deterministic generation when you want different values
// but don't want to fill a slice.
func (a *Autofill) FillWithIndex(v interface{}, index int, overrides ...Override) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("Fill requires a pointer to struct, got %T", v)
	}

	elem := rv.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("Fill requires a pointer to struct, got pointer to %s", elem.Kind())
	}

	// Merge overrides
	override := mergeOverrides(overrides)

	// Create context
	ctx := newContext(a.locale, a.seed, index, a.rand)
	ctx = ctx.withStruct(v)

	// Fill each field
	typ := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := elem.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		fieldCtx := ctx.withFieldName(field.Name)

		// Check for override
		if overrideVal, ok := override[field.Name]; ok {
			resolved := resolveOverride(overrideVal, index)
			if resolved != nil {
				if err := setFieldValue(fieldVal, resolved); err != nil {
					return fmt.Errorf("failed to set override for field %s: %w", field.Name, err)
				}
				continue
			}
		}

		// Generate value
		val, err := a.generateValue(field, fieldCtx)
		if err != nil {
			return fmt.Errorf("failed to generate value for field %s: %w", field.Name, err)
		}

		if val != nil {
			if err := setFieldValue(fieldVal, val); err != nil {
				return fmt.Errorf("failed to set value for field %s: %w", field.Name, err)
			}
		}
	}

	return nil
}

// FillSlice populates a slice of structs with generated test data.
// The input must be a pointer to a slice of structs.
// Optional overrides can be provided to set specific field values.
// Overrides can include SequenceFunc values for index-based generation.
//
// Example:
//
//	users := make([]User, 10)
//	err := autofill.New().FillSlice(&users, autofill.Override{
//	    "Email": autofill.Seq("user%d@example.com"),
//	})
func (a *Autofill) FillSlice(v interface{}, overrides ...Override) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("FillSlice requires a pointer to slice, got %T", v)
	}

	elem := rv.Elem()
	if elem.Kind() != reflect.Slice {
		return fmt.Errorf("FillSlice requires a pointer to slice, got pointer to %s", elem.Kind())
	}

	sliceLen := elem.Len()
	for i := 0; i < sliceLen; i++ {
		item := elem.Index(i).Addr().Interface()
		if err := a.FillWithIndex(item, i, overrides...); err != nil {
			return fmt.Errorf("failed to fill slice element at index %d: %w", i, err)
		}
	}

	return nil
}

// Fill is a convenience function that creates a new Autofill instance
// and fills the given struct.
func Fill(v interface{}, overrides ...Override) error {
	return New().Fill(v, overrides...)
}

// FillSlice is a convenience function that creates a new Autofill instance
// and fills the given slice of structs.
func FillSlice(v interface{}, overrides ...Override) error {
	return New().FillSlice(v, overrides...)
}

// mergeOverrides merges multiple Override maps into one.
// Later overrides take precedence over earlier ones.
func mergeOverrides(overrides []Override) Override {
	if len(overrides) == 0 {
		return nil
	}
	if len(overrides) == 1 {
		return overrides[0]
	}

	merged := make(Override)
	for _, override := range overrides {
		for k, v := range override {
			merged[k] = v
		}
	}
	return merged
}

// setFieldValue sets a reflect.Value with the given value, handling type conversions.
func setFieldValue(field reflect.Value, value interface{}) error {
	if value == nil {
		return nil
	}

	valReflect := reflect.ValueOf(value)
	fieldType := field.Type()

	// Handle pointer types
	if fieldType.Kind() == reflect.Ptr {
		if valReflect.Type().ConvertibleTo(fieldType.Elem()) {
			ptr := reflect.New(fieldType.Elem())
			ptr.Elem().Set(valReflect.Convert(fieldType.Elem()))
			field.Set(ptr)
			return nil
		}
		if valReflect.Type() == fieldType {
			field.Set(valReflect)
			return nil
		}
	}

	// Direct assignment if types match
	if valReflect.Type() == fieldType {
		field.Set(valReflect)
		return nil
	}

	// Try conversion
	if valReflect.Type().ConvertibleTo(fieldType) {
		field.Set(valReflect.Convert(fieldType))
		return nil
	}

	return fmt.Errorf("cannot set field of type %s with value of type %T", fieldType, value)
}
