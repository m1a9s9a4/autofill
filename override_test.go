package autofill

import "testing"

func TestSeq(t *testing.T) {
	fn := Seq("user%d@example.com")

	tests := []struct {
		index    int
		expected string
	}{
		{0, "user0@example.com"},
		{1, "user1@example.com"},
		{5, "user5@example.com"},
		{100, "user100@example.com"},
	}

	for _, tt := range tests {
		result := fn(tt.index)
		if result != tt.expected {
			t.Errorf("Seq(%d) = %v, expected %v", tt.index, result, tt.expected)
		}
	}
}

func TestSeqInt(t *testing.T) {
	fn := SeqInt(100)

	tests := []struct {
		index    int
		expected int
	}{
		{0, 100},
		{1, 101},
		{5, 105},
		{100, 200},
	}

	for _, tt := range tests {
		result := fn(tt.index)
		if result != tt.expected {
			t.Errorf("SeqInt(%d) = %v, expected %v", tt.index, result, tt.expected)
		}
	}
}

func TestSeqInt64(t *testing.T) {
	fn := SeqInt64(1000)

	tests := []struct {
		index    int
		expected int64
	}{
		{0, 1000},
		{1, 1001},
		{5, 1005},
		{100, 1100},
	}

	for _, tt := range tests {
		result := fn(tt.index)
		if result != tt.expected {
			t.Errorf("SeqInt64(%d) = %v, expected %v", tt.index, result, tt.expected)
		}
	}
}

func TestRandom(t *testing.T) {
	fn := Random(1, 10)

	for i := 0; i < 100; i++ {
		result := fn(i)
		n, ok := result.(int)
		if !ok {
			t.Fatalf("expected int, got %T", result)
		}
		if n < 1 || n > 10 {
			t.Errorf("Random(%d) = %d, expected value in range [1, 10]", i, n)
		}
	}
}

func TestResolveOverride_WithSequenceFunc(t *testing.T) {
	fn := Seq("test%d")
	result := resolveOverride(fn, 5)

	expected := "test5"
	if result != expected {
		t.Errorf("resolveOverride(SeqFunc, 5) = %v, expected %v", result, expected)
	}
}

func TestResolveOverride_WithDirectValue(t *testing.T) {
	value := "direct value"
	result := resolveOverride(value, 5)

	if result != value {
		t.Errorf("resolveOverride(direct, 5) = %v, expected %v", result, value)
	}
}

func TestResolveOverride_WithInt(t *testing.T) {
	value := 42
	result := resolveOverride(value, 5)

	if result != value {
		t.Errorf("resolveOverride(int, 5) = %v, expected %v", result, value)
	}
}

func TestOverride_Integration(t *testing.T) {
	type TestStruct struct {
		Name  string
		Email string
		Age   int
		ID    int64
	}

	tests := []TestStruct{
		{},
		{},
		{},
	}

	err := FillSlice(&tests, Override{
		"Email": Seq("user%d@example.com"),
		"Age":   SeqInt(20),
		"ID":    SeqInt64(1000),
		"Name":  "Fixed Name",
	})

	if err != nil {
		t.Fatalf("FillSlice failed: %v", err)
	}

	expected := []struct {
		Email string
		Age   int
		ID    int64
		Name  string
	}{
		{"user0@example.com", 20, 1000, "Fixed Name"},
		{"user1@example.com", 21, 1001, "Fixed Name"},
		{"user2@example.com", 22, 1002, "Fixed Name"},
	}

	for i, exp := range expected {
		if tests[i].Email != exp.Email {
			t.Errorf("test[%d].Email = %s, expected %s", i, tests[i].Email, exp.Email)
		}
		if tests[i].Age != exp.Age {
			t.Errorf("test[%d].Age = %d, expected %d", i, tests[i].Age, exp.Age)
		}
		if tests[i].ID != exp.ID {
			t.Errorf("test[%d].ID = %d, expected %d", i, tests[i].ID, exp.ID)
		}
		if tests[i].Name != exp.Name {
			t.Errorf("test[%d].Name = %s, expected %s", i, tests[i].Name, exp.Name)
		}
	}
}
