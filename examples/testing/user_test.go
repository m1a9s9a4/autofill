package testing

import (
	"fmt"
	"testing"

	"github.com/m1a9s9a4/autofill"
)

type User struct {
	ID       int64  `autofill:"seq"`
	Name     string
	Email    string `autofill:"email"`
	Age      int    `autofill:"min=18,max=65"`
	Username string
}

func TestUserCreation(t *testing.T) {
	var user User
	err := autofill.Fill(&user)

	if err != nil {
		t.Fatalf("Fill failed: %v", err)
	}

	// ID starts at 0 with seq tag - this is expected
	if user.Name == "" {
		t.Error("Name should not be empty")
	}
	if user.Email == "" {
		t.Error("Email should not be empty")
	}
	if user.Age < 18 || user.Age > 65 {
		t.Errorf("Age %d is out of range [18, 65]", user.Age)
	}
}

func TestUserBulkCreation(t *testing.T) {
	users := make([]User, 100)
	err := autofill.FillSlice(&users)

	if err != nil {
		t.Fatalf("FillSlice failed: %v", err)
	}

	// Check all users have unique IDs (due to seq tag)
	ids := make(map[int64]bool)
	for i, user := range users {
		if user.Email == "" {
			t.Errorf("user %d: Email should not be empty", i)
		}
		if ids[user.ID] {
			t.Errorf("user %d: Duplicate ID %d", i, user.ID)
		}
		ids[user.ID] = true
	}
}

func TestUserWithOverride(t *testing.T) {
	users := make([]User, 5)
	err := autofill.FillSlice(&users, autofill.Override{
		"Email":    autofill.Seq("testuser%d@example.com"),
		"Username": autofill.Seq("user%d"),
		"Age":      autofill.SeqInt(20),
	})

	if err != nil {
		t.Fatalf("FillSlice with override failed: %v", err)
	}

	for i, user := range users {
		expectedEmail := fmt.Sprintf("testuser%d@example.com", i)
		if user.Email != expectedEmail {
			t.Errorf("user %d: expected Email %s, got %s", i, expectedEmail, user.Email)
		}

		expectedUsername := fmt.Sprintf("user%d", i)
		if user.Username != expectedUsername {
			t.Errorf("user %d: expected Username %s, got %s", i, expectedUsername, user.Username)
		}

		expectedAge := 20 + i
		if user.Age != expectedAge {
			t.Errorf("user %d: expected Age %d, got %d", i, expectedAge, user.Age)
		}
	}
}

func TestDeterministicGeneration(t *testing.T) {
	af := autofill.New().WithSeed(42)

	var user1 User
	if err := af.Fill(&user1); err != nil {
		t.Fatalf("Fill failed: %v", err)
	}

	af2 := autofill.New().WithSeed(42)
	var user2 User
	if err := af2.Fill(&user2); err != nil {
		t.Fatalf("Fill failed: %v", err)
	}

	// With same seed, all non-random fields should be the same
	if user1.ID != user2.ID {
		t.Errorf("expected same ID, got %d and %d", user1.ID, user2.ID)
	}
	if user1.Name != user2.Name {
		t.Errorf("expected same Name, got %s and %s", user1.Name, user2.Name)
	}
	if user1.Age != user2.Age {
		t.Errorf("expected same Age, got %d and %d", user1.Age, user2.Age)
	}
}

func BenchmarkFillUser(b *testing.B) {
	af := autofill.New()
	var user User

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		af.Fill(&user)
	}
}

func BenchmarkFillUserSlice(b *testing.B) {
	af := autofill.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		users := make([]User, 100)
		af.FillSlice(&users)
	}
}
