package autofill

import (
	"testing"
	"time"
)

type TestUser struct {
	ID        int64
	Name      string
	Email     string
	Age       int
	Active    bool
	Score     float64
	CreatedAt time.Time
}

func TestNew(t *testing.T) {
	af := New()
	if af == nil {
		t.Fatal("New() returned nil")
	}
	if af.locale != "en_US" {
		t.Errorf("expected default locale en_US, got %s", af.locale)
	}
	if af.seed == 0 {
		t.Error("expected non-zero seed")
	}
	if af.rules == nil {
		t.Error("expected rules to be initialized")
	}
}

func TestWithLocale(t *testing.T) {
	af := New().WithLocale("ja_JP")
	if af.locale != "ja_JP" {
		t.Errorf("expected locale ja_JP, got %s", af.locale)
	}
}

func TestWithSeed(t *testing.T) {
	af := New().WithSeed(12345)
	if af.seed != 12345 {
		t.Errorf("expected seed 12345, got %d", af.seed)
	}
}

func TestFill_BasicTypes(t *testing.T) {
	var user TestUser
	err := Fill(&user)
	if err != nil {
		t.Fatalf("Fill failed: %v", err)
	}

	if user.Name == "" {
		t.Error("Name should not be empty")
	}
	if user.Email == "" {
		t.Error("Email should not be empty")
	}
	if user.Age == 0 {
		t.Error("Age should not be zero")
	}
	if user.Score == 0 {
		t.Error("Score should not be zero")
	}
	if user.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestFill_WithOverride(t *testing.T) {
	var user TestUser
	err := Fill(&user, Override{
		"Name":  "John Doe",
		"Email": "john@example.com",
		"Age":   30,
	})
	if err != nil {
		t.Fatalf("Fill with override failed: %v", err)
	}

	if user.Name != "John Doe" {
		t.Errorf("expected Name to be John Doe, got %s", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Errorf("expected Email to be john@example.com, got %s", user.Email)
	}
	if user.Age != 30 {
		t.Errorf("expected Age to be 30, got %d", user.Age)
	}
}

func TestFill_Deterministic(t *testing.T) {
	af := New().WithSeed(12345)

	var user1 TestUser
	if err := af.Fill(&user1); err != nil {
		t.Fatalf("Fill failed: %v", err)
	}

	af2 := New().WithSeed(12345)
	var user2 TestUser
	if err := af2.Fill(&user2); err != nil {
		t.Fatalf("Fill failed: %v", err)
	}

	if user1.Name != user2.Name {
		t.Errorf("expected same Name, got %s and %s", user1.Name, user2.Name)
	}
	if user1.Age != user2.Age {
		t.Errorf("expected same Age, got %d and %d", user1.Age, user2.Age)
	}
	if user1.Active != user2.Active {
		t.Errorf("expected same Active, got %v and %v", user1.Active, user2.Active)
	}
}

func TestFillSlice(t *testing.T) {
	users := make([]TestUser, 5)
	err := FillSlice(&users)
	if err != nil {
		t.Fatalf("FillSlice failed: %v", err)
	}

	for i, user := range users {
		if user.Name == "" {
			t.Errorf("user %d: Name should not be empty", i)
		}
		if user.Email == "" {
			t.Errorf("user %d: Email should not be empty", i)
		}
	}

	// Check that values are different
	if users[0].Name == users[1].Name && users[1].Name == users[2].Name {
		t.Error("expected different names for different users")
	}
}

func TestFillSlice_WithSequence(t *testing.T) {
	users := make([]TestUser, 3)
	err := FillSlice(&users, Override{
		"Email": Seq("user%d@example.com"),
		"Age":   SeqInt(20),
	})
	if err != nil {
		t.Fatalf("FillSlice with sequence failed: %v", err)
	}

	expected := []struct {
		Email string
		Age   int
	}{
		{"user0@example.com", 20},
		{"user1@example.com", 21},
		{"user2@example.com", 22},
	}

	for i, exp := range expected {
		if users[i].Email != exp.Email {
			t.Errorf("user %d: expected Email %s, got %s", i, exp.Email, users[i].Email)
		}
		if users[i].Age != exp.Age {
			t.Errorf("user %d: expected Age %d, got %d", i, exp.Age, users[i].Age)
		}
	}
}

func TestFill_WithTags(t *testing.T) {
	type TaggedUser struct {
		ID     int64  `autofill:"seq"`
		Email  string `autofill:"email"`
		URL    string `autofill:"url"`
		Status string `autofill:"oneof=active|inactive"`
		Skip   string `autofill:"-"`
	}

	var user TaggedUser
	err := Fill(&user)
	if err != nil {
		t.Fatalf("Fill with tags failed: %v", err)
	}

	if user.ID != 0 {
		t.Logf("ID: %d", user.ID)
	}
	if user.Email == "" {
		t.Error("Email should not be empty")
	}
	if user.URL == "" {
		t.Error("URL should not be empty")
	}
	if user.Status != "active" && user.Status != "inactive" {
		t.Errorf("Status should be active or inactive, got %s", user.Status)
	}
	if user.Skip != "" {
		t.Errorf("Skip should be empty, got %s", user.Skip)
	}
}

func TestFill_WithMinMax(t *testing.T) {
	type RangedStruct struct {
		Age int `autofill:"min=18,max=30"`
	}

	for i := 0; i < 20; i++ {
		var s RangedStruct
		err := New().FillWithIndex(&s, i)
		if err != nil {
			t.Fatalf("Fill failed: %v", err)
		}
		if s.Age < 18 || s.Age > 30 {
			t.Errorf("Age %d is out of range [18, 30]", s.Age)
		}
	}
}

func TestFill_NestedStruct(t *testing.T) {
	type Address struct {
		City    string
		Country string
	}

	type Person struct {
		Name    string
		Address Address
	}

	var person Person
	err := Fill(&person)
	if err != nil {
		t.Fatalf("Fill nested struct failed: %v", err)
	}

	if person.Name == "" {
		t.Error("Name should not be empty")
	}
	if person.Address.City == "" {
		t.Error("Address.City should not be empty")
	}
	if person.Address.Country == "" {
		t.Error("Address.Country should not be empty")
	}
}

func TestFill_PointerFields(t *testing.T) {
	type PointerStruct struct {
		Name  *string
		Age   *int
		Email *string
	}

	var s PointerStruct
	err := Fill(&s)
	if err != nil {
		t.Fatalf("Fill with pointer fields failed: %v", err)
	}

	if s.Name == nil {
		t.Error("Name pointer should not be nil")
	} else if *s.Name == "" {
		t.Error("Name value should not be empty")
	}

	if s.Age == nil {
		t.Error("Age pointer should not be nil")
	} else if *s.Age == 0 {
		t.Error("Age value should not be zero")
	}
}

func TestFill_InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"not a pointer", TestUser{}},
		{"pointer to non-struct", new(int)},
		{"nil", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Fill(tt.input)
			if err == nil {
				t.Error("expected error for invalid input, got nil")
			}
		})
	}
}

func TestFillSlice_InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"not a pointer", []TestUser{}},
		{"pointer to non-slice", new(TestUser)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FillSlice(tt.input)
			if err == nil {
				t.Error("expected error for invalid input, got nil")
			}
		})
	}
}

func BenchmarkFill(b *testing.B) {
	af := New()
	var user TestUser

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		af.Fill(&user)
	}
}

func BenchmarkFillSlice(b *testing.B) {
	af := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		users := make([]TestUser, 100)
		af.FillSlice(&users)
	}
}
