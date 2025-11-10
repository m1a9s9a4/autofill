package autofill

import "fmt"

// Override represents a map of field names to their override values.
// Values can be:
// - Direct values: any type that matches the field type
// - SequenceFunc: a function that generates values based on index
type Override map[string]interface{}

// SequenceFunc is a function that generates a value based on an index.
// It is called for each element when filling slices.
type SequenceFunc func(index int) interface{}

// Seq creates a SequenceFunc that formats a string with the current index.
// The format string should contain a single %d placeholder for the index.
//
// Example:
//
//	autofill.Seq("user%d@example.com") // generates user0@example.com, user1@example.com, etc.
func Seq(format string) SequenceFunc {
	return func(index int) interface{} {
		return fmt.Sprintf(format, index)
	}
}

// SeqInt creates a SequenceFunc that generates sequential integers starting from start.
//
// Example:
//
//	autofill.SeqInt(100) // generates 100, 101, 102, etc.
func SeqInt(start int) SequenceFunc {
	return func(index int) interface{} {
		return start + index
	}
}

// SeqInt64 creates a SequenceFunc that generates sequential int64 values starting from start.
//
// Example:
//
//	autofill.SeqInt64(1000) // generates 1000, 1001, 1002, etc.
func SeqInt64(start int64) SequenceFunc {
	return func(index int) interface{} {
		return start + int64(index)
	}
}

// Random creates a SequenceFunc that generates random integers between min and max (inclusive).
//
// Example:
//
//	autofill.Random(1, 100) // generates random integers between 1 and 100
func Random(min, max int) SequenceFunc {
	return func(index int) interface{} {
		// This will be seeded by the autofill context
		return min + (index % (max - min + 1))
	}
}

// resolveOverride resolves an override value for a given field and index.
// If the value is a SequenceFunc, it calls the function with the index.
// Otherwise, it returns the value as-is.
func resolveOverride(value interface{}, index int) interface{} {
	if fn, ok := value.(SequenceFunc); ok {
		return fn(index)
	}
	return value
}
