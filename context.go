package autofill

import (
	"math/rand"
	"reflect"
)

// Context provides information about the current generation context.
// It is passed to Rule.Generate() methods to allow context-aware value generation.
type Context interface {
	// Locale returns the current locale (e.g., "ja_JP", "en_US")
	Locale() string

	// Seed returns the random seed used for generation
	Seed() int64

	// Index returns the current index when filling slices (0-based)
	Index() int

	// Rand returns the random number generator for this context
	Rand() *rand.Rand

	// GetField returns the value of a field by name from the current struct being filled
	GetField(name string) (interface{}, bool)

	// GetStruct returns the struct being filled
	GetStruct() interface{}

	// FieldName returns the name of the current field being filled
	FieldName() string
}

// context is the internal implementation of Context
type context struct {
	locale    string
	seed      int64
	index     int
	rand      *rand.Rand
	fieldMap  map[string]interface{}
	structVal interface{}
	fieldName string
}

// NewContext creates a new Context with the given parameters
func newContext(locale string, seed int64, index int, r *rand.Rand) *context {
	return &context{
		locale:   locale,
		seed:     seed,
		index:    index,
		rand:     r,
		fieldMap: make(map[string]interface{}),
	}
}

// Locale returns the current locale
func (c *context) Locale() string {
	return c.locale
}

// Seed returns the random seed
func (c *context) Seed() int64 {
	return c.seed
}

// Index returns the current index
func (c *context) Index() int {
	return c.index
}

// Rand returns the random number generator
func (c *context) Rand() *rand.Rand {
	return c.rand
}

// GetField returns the value of a field by name
func (c *context) GetField(name string) (interface{}, bool) {
	val, ok := c.fieldMap[name]
	return val, ok
}

// GetStruct returns the struct being filled
func (c *context) GetStruct() interface{} {
	return c.structVal
}

// FieldName returns the name of the current field
func (c *context) FieldName() string {
	return c.fieldName
}

// withStruct creates a new context with the given struct value
func (c *context) withStruct(v interface{}) *context {
	newCtx := *c
	newCtx.structVal = v
	newCtx.fieldMap = make(map[string]interface{})

	// Populate field map
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() == reflect.Struct {
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if field.IsValid() && field.CanInterface() {
				newCtx.fieldMap[typ.Field(i).Name] = field.Interface()
			}
		}
	}

	return &newCtx
}

// withFieldName creates a new context with the given field name
func (c *context) withFieldName(name string) *context {
	newCtx := *c
	newCtx.fieldName = name
	return &newCtx
}

// withIndex creates a new context with the given index
func (c *context) withIndex(index int) *context {
	newCtx := *c
	newCtx.index = index
	return &newCtx
}
