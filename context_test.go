package autofill

import (
	"math/rand"
	"testing"
)

func TestNewContext(t *testing.T) {
	r := rand.New(rand.NewSource(12345))
	ctx := newContext("ja_JP", 12345, 0, r)

	if ctx == nil {
		t.Fatal("newContext returned nil")
	}

	if ctx.Locale() != "ja_JP" {
		t.Errorf("expected locale ja_JP, got %s", ctx.Locale())
	}

	if ctx.Seed() != 12345 {
		t.Errorf("expected seed 12345, got %d", ctx.Seed())
	}

	if ctx.Index() != 0 {
		t.Errorf("expected index 0, got %d", ctx.Index())
	}

	if ctx.Rand() != r {
		t.Error("expected same rand instance")
	}
}

func TestContext_GetField(t *testing.T) {
	type TestStruct struct {
		Name string
		Age  int
	}

	s := TestStruct{Name: "John", Age: 30}
	r := rand.New(rand.NewSource(12345))
	ctx := newContext("en_US", 12345, 0, r)
	ctx = ctx.withStruct(&s)

	val, ok := ctx.GetField("Name")
	if !ok {
		t.Error("expected to find Name field")
	}
	if val != "John" {
		t.Errorf("expected Name to be John, got %v", val)
	}

	val, ok = ctx.GetField("Age")
	if !ok {
		t.Error("expected to find Age field")
	}
	if val != 30 {
		t.Errorf("expected Age to be 30, got %v", val)
	}

	_, ok = ctx.GetField("NonExistent")
	if ok {
		t.Error("expected not to find NonExistent field")
	}
}

func TestContext_GetStruct(t *testing.T) {
	type TestStruct struct {
		Name string
	}

	s := TestStruct{Name: "John"}
	r := rand.New(rand.NewSource(12345))
	ctx := newContext("en_US", 12345, 0, r)
	ctx = ctx.withStruct(&s)

	structVal := ctx.GetStruct()
	if structVal == nil {
		t.Error("expected non-nil struct")
	}

	// Verify we can access the struct
	ptr, ok := structVal.(*TestStruct)
	if !ok {
		t.Fatalf("expected *TestStruct, got %T", structVal)
	}
	if ptr.Name != "John" {
		t.Errorf("expected Name to be John, got %s", ptr.Name)
	}
}

func TestContext_FieldName(t *testing.T) {
	r := rand.New(rand.NewSource(12345))
	ctx := newContext("en_US", 12345, 0, r)
	ctx = ctx.withFieldName("TestField")

	if ctx.FieldName() != "TestField" {
		t.Errorf("expected FieldName to be TestField, got %s", ctx.FieldName())
	}
}

func TestContext_WithIndex(t *testing.T) {
	r := rand.New(rand.NewSource(12345))
	ctx := newContext("en_US", 12345, 0, r)

	newCtx := ctx.withIndex(5)

	if newCtx.Index() != 5 {
		t.Errorf("expected index 5, got %d", newCtx.Index())
	}

	// Original context should be unchanged
	if ctx.Index() != 0 {
		t.Errorf("original context index should still be 0, got %d", ctx.Index())
	}
}

func TestContext_WithStruct(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 int
	}

	s := TestStruct{Field1: "test", Field2: 42}
	r := rand.New(rand.NewSource(12345))
	ctx := newContext("en_US", 12345, 0, r)

	newCtx := ctx.withStruct(&s)

	if newCtx.GetStruct() == nil {
		t.Error("expected struct to be set")
	}

	val, ok := newCtx.GetField("Field1")
	if !ok || val != "test" {
		t.Errorf("expected Field1 to be 'test', got %v (ok=%v)", val, ok)
	}

	val, ok = newCtx.GetField("Field2")
	if !ok || val != 42 {
		t.Errorf("expected Field2 to be 42, got %v (ok=%v)", val, ok)
	}
}

func TestContext_WithFieldName(t *testing.T) {
	r := rand.New(rand.NewSource(12345))
	ctx := newContext("en_US", 12345, 0, r)

	newCtx := ctx.withFieldName("NewField")

	if newCtx.FieldName() != "NewField" {
		t.Errorf("expected field name NewField, got %s", newCtx.FieldName())
	}

	// Original context should be unchanged
	if ctx.FieldName() != "" {
		t.Errorf("original context field name should be empty, got %s", ctx.FieldName())
	}
}

func TestContext_ImmutabilityChain(t *testing.T) {
	r := rand.New(rand.NewSource(12345))
	ctx1 := newContext("en_US", 12345, 0, r)
	ctx2 := ctx1.withIndex(1)
	ctx3 := ctx2.withFieldName("Field1")

	// All contexts should have different values
	if ctx1.Index() != 0 {
		t.Errorf("ctx1 index should be 0, got %d", ctx1.Index())
	}
	if ctx2.Index() != 1 {
		t.Errorf("ctx2 index should be 1, got %d", ctx2.Index())
	}
	if ctx3.Index() != 1 {
		t.Errorf("ctx3 index should be 1, got %d", ctx3.Index())
	}

	if ctx1.FieldName() != "" {
		t.Errorf("ctx1 field name should be empty, got %s", ctx1.FieldName())
	}
	if ctx2.FieldName() != "" {
		t.Errorf("ctx2 field name should be empty, got %s", ctx2.FieldName())
	}
	if ctx3.FieldName() != "Field1" {
		t.Errorf("ctx3 field name should be Field1, got %s", ctx3.FieldName())
	}
}
