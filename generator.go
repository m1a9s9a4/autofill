package autofill

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

// generateValue generates a value for the given field based on its type and tags.
func (a *Autofill) generateValue(field reflect.StructField, ctx *context) (interface{}, error) {
	// Check for autofill tag
	tag := field.Tag.Get("autofill")
	if tag == "-" {
		return nil, nil // Skip this field
	}

	// Parse tag if present
	if tag != "" {
		val, err := a.generateFromTag(tag, field, ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to generate from tag %q: %w", tag, err)
		}
		if val != nil {
			return val, nil
		}
	}

	// Generate based on type
	return a.generateByType(field.Type, ctx)
}

// generateFromTag generates a value based on the autofill struct tag.
func (a *Autofill) generateFromTag(tag string, field reflect.StructField, ctx *context) (interface{}, error) {
	parts := strings.Split(tag, ",")

	// Parse all tag parts into params
	params := parseTagParams(parts)

	// Check for rule=<name> format
	if ruleName, ok := params["rule"]; ok {
		if a.rules != nil {
			if rule, ok := a.rules.Get(ruleName); ok {
				return rule.Generate(ctx)
			}
		}
		return nil, fmt.Errorf("rule %q not found", ruleName)
	}

	// Handle built-in tags (first part without =)
	mainTag := strings.TrimSpace(parts[0])
	if !strings.Contains(mainTag, "=") {
		switch mainTag {
		case "seq":
			return int64(ctx.Index()), nil
		case "now":
			return time.Now(), nil
		case "email":
			return a.generateEmail(ctx), nil
		case "url":
			return a.generateURL(ctx), nil
		case "uuid":
			return a.generateUUID(ctx), nil
		}

		// Check if it's a direct rule name (without "rule=" prefix)
		if a.rules != nil {
			if rule, ok := a.rules.Get(mainTag); ok {
				return rule.Generate(ctx)
			}
		}
	}

	// Handle min/max for numeric types
	if minStr, hasMin := params["min"]; hasMin {
		if maxStr, hasMax := params["max"]; hasMax {
			var min, max int
			fmt.Sscanf(minStr, "%d", &min)
			fmt.Sscanf(maxStr, "%d", &max)
			if min <= max {
				rangeSize := max - min + 1
				return min + (ctx.Index() % rangeSize), nil
			}
		}
	}

	// Handle oneof
	if oneofStr, ok := params["oneof"]; ok {
		options := strings.Split(oneofStr, "|")
		if len(options) > 0 {
			return options[ctx.Index()%len(options)], nil
		}
	}

	return nil, nil
}

// parseTagParams parses tag parameters in key=value format.
func parseTagParams(parts []string) map[string]string {
	params := make(map[string]string)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if kv := strings.SplitN(part, "=", 2); len(kv) == 2 {
			params[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return params
}

// generateByType generates a value based on the reflect.Type.
func (a *Autofill) generateByType(typ reflect.Type, ctx *context) (interface{}, error) {
	switch typ.Kind() {
	case reflect.String:
		return a.generateString(ctx), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.generateInt(ctx), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return a.generateUint(ctx), nil
	case reflect.Float32, reflect.Float64:
		return a.generateFloat(ctx), nil
	case reflect.Bool:
		return a.generateBool(ctx), nil
	case reflect.Ptr:
		// Generate value for the element type
		val, err := a.generateByType(typ.Elem(), ctx)
		if err != nil {
			return nil, err
		}
		// Create pointer to the value
		ptr := reflect.New(typ.Elem())
		ptr.Elem().Set(reflect.ValueOf(val).Convert(typ.Elem()))
		return ptr.Interface(), nil
	case reflect.Slice:
		// Generate a slice with 3 elements by default
		return a.generateSlice(typ, ctx, 3)
	case reflect.Struct:
		// Special handling for time.Time
		if typ == reflect.TypeOf(time.Time{}) {
			return a.generateTime(ctx), nil
		}
		// For other structs, recursively fill
		return a.fillStructValue(typ, ctx)
	default:
		return nil, fmt.Errorf("unsupported type: %s", typ)
	}
}

// generateString generates a random string.
func (a *Autofill) generateString(ctx *context) string {
	words := []string{"hello", "world", "test", "sample", "data", "value", "string", "text"}
	return words[ctx.Index()%len(words)]
}

// generateInt generates a random integer.
func (a *Autofill) generateInt(ctx *context) int64 {
	return int64(100 + (ctx.Index() % 900))
}

// generateUint generates a random unsigned integer.
func (a *Autofill) generateUint(ctx *context) uint64 {
	return uint64(100 + (ctx.Index() % 900))
}

// generateFloat generates a random float.
func (a *Autofill) generateFloat(ctx *context) float64 {
	seed := ctx.Seed() + int64(ctx.Index())
	rng := rand.New(rand.NewSource(seed))
	return rng.Float64() * 1000
}

// generateBool generates a random boolean.
func (a *Autofill) generateBool(ctx *context) bool {
	return ctx.Index()%2 == 0
}

// generateTime generates a random time.
func (a *Autofill) generateTime(ctx *context) time.Time {
	// Generate a time within the last year
	now := time.Now()
	daysAgo := ctx.Index() % 365
	return now.AddDate(0, 0, -daysAgo)
}

// generateEmail generates an email address.
func (a *Autofill) generateEmail(ctx *context) string {
	domains := []string{"example.com", "test.com", "mail.com"}
	prefixes := []string{"user", "test", "demo", "sample"}

	prefix := prefixes[ctx.Index()%len(prefixes)]
	domain := domains[ctx.Index()%len(domains)]

	return fmt.Sprintf("%s%d@%s", prefix, ctx.Index(), domain)
}

// generateURL generates a URL.
func (a *Autofill) generateURL(ctx *context) string {
	domains := []string{"example.com", "test.com", "demo.org"}
	paths := []string{"/", "/home", "/about", "/contact"}

	domain := domains[ctx.Index()%len(domains)]
	path := paths[ctx.Index()%len(paths)]

	return fmt.Sprintf("https://%s%s", domain, path)
}

// generateUUID generates a UUID.
func (a *Autofill) generateUUID(ctx *context) string {
	seed := ctx.Seed() + int64(ctx.Index())
	rng := rand.New(rand.NewSource(seed))

	var uuidBytes [16]byte
	rng.Read(uuidBytes[:])

	// Set version (4) and variant bits
	uuidBytes[6] = (uuidBytes[6] & 0x0f) | 0x40
	uuidBytes[8] = (uuidBytes[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuidBytes[0:4],
		uuidBytes[4:6],
		uuidBytes[6:8],
		uuidBytes[8:10],
		uuidBytes[10:16])
}

// generateSlice generates a slice of the given type and length.
func (a *Autofill) generateSlice(typ reflect.Type, ctx *context, length int) (interface{}, error) {
	slice := reflect.MakeSlice(typ, length, length)
	elemType := typ.Elem()

	for i := 0; i < length; i++ {
		elemCtx := ctx.withIndex(i)
		val, err := a.generateByType(elemType, elemCtx)
		if err != nil {
			return nil, fmt.Errorf("failed to generate slice element at index %d: %w", i, err)
		}
		slice.Index(i).Set(reflect.ValueOf(val).Convert(elemType))
	}

	return slice.Interface(), nil
}

// fillStructValue fills a struct value and returns it.
func (a *Autofill) fillStructValue(typ reflect.Type, ctx *context) (interface{}, error) {
	structVal := reflect.New(typ).Elem()
	structCtx := ctx.withStruct(structVal.Addr().Interface())

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := structVal.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		fieldCtx := structCtx.withFieldName(field.Name)
		val, err := a.generateValue(field, fieldCtx)
		if err != nil {
			return nil, fmt.Errorf("failed to generate field %s: %w", field.Name, err)
		}

		if val != nil {
			fieldVal.Set(reflect.ValueOf(val).Convert(field.Type))
		}
	}

	return structVal.Interface(), nil
}
