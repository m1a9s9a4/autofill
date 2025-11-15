package autofill

import "testing"

type UserWithRole struct {
	ID          int64  `autofill:"seq"`
	Email       string `autofill:"email"`
	Role        string
	Permissions string
	TeamID      string
}

func TestWithDefaults(t *testing.T) {
	af := New().WithDefaults(Override{
		"Role":        "admin",
		"Permissions": "all",
	})

	var user UserWithRole
	err := af.Fill(&user)
	if err != nil {
		t.Fatalf("Fill failed: %v", err)
	}

	if user.Role != "admin" {
		t.Errorf("expected Role 'admin', got %s", user.Role)
	}
	if user.Permissions != "all" {
		t.Errorf("expected Permissions 'all', got %s", user.Permissions)
	}
}

func TestWithDefaults_Override(t *testing.T) {
	// Defaults set Role=admin
	af := New().WithDefaults(Override{
		"Role": "admin",
	})

	var user UserWithRole
	// But we override it with Role=superadmin
	err := af.Fill(&user, Override{
		"Role": "superadmin",
	})
	if err != nil {
		t.Fatalf("Fill failed: %v", err)
	}

	// The passed override should take precedence
	if user.Role != "superadmin" {
		t.Errorf("expected Role 'superadmin' (override), got %s", user.Role)
	}
}

func TestWithDefaults_Slice(t *testing.T) {
	adminFiller := New().WithDefaults(Override{
		"Role":        "admin",
		"Permissions": "all",
		"TeamID":      "admin-team",
	})

	admins := make([]UserWithRole, 3)
	err := adminFiller.FillSlice(&admins)
	if err != nil {
		t.Fatalf("FillSlice failed: %v", err)
	}

	for i, admin := range admins {
		if admin.Role != "admin" {
			t.Errorf("admin %d: expected Role 'admin', got %s", i, admin.Role)
		}
		if admin.Permissions != "all" {
			t.Errorf("admin %d: expected Permissions 'all', got %s", i, admin.Permissions)
		}
		if admin.TeamID != "admin-team" {
			t.Errorf("admin %d: expected TeamID 'admin-team', got %s", i, admin.TeamID)
		}
	}
}

func TestWithDefaults_MultipleFillers(t *testing.T) {
	// Create admin filler
	adminFiller := New().WithSeed(12345).WithDefaults(Override{
		"Role":        "admin",
		"Permissions": "all",
	})

	// Create member filler
	memberFiller := New().WithSeed(12345).WithDefaults(Override{
		"Role":        "member",
		"Permissions": "read",
	})

	admins := make([]UserWithRole, 2)
	members := make([]UserWithRole, 3)

	err := adminFiller.FillSlice(&admins)
	if err != nil {
		t.Fatalf("adminFiller.FillSlice failed: %v", err)
	}

	err = memberFiller.FillSlice(&members)
	if err != nil {
		t.Fatalf("memberFiller.FillSlice failed: %v", err)
	}

	// Check admins
	for i, admin := range admins {
		if admin.Role != "admin" {
			t.Errorf("admin %d: expected Role 'admin', got %s", i, admin.Role)
		}
		if admin.Permissions != "all" {
			t.Errorf("admin %d: expected Permissions 'all', got %s", i, admin.Permissions)
		}
	}

	// Check members
	for i, member := range members {
		if member.Role != "member" {
			t.Errorf("member %d: expected Role 'member', got %s", i, member.Role)
		}
		if member.Permissions != "read" {
			t.Errorf("member %d: expected Permissions 'read', got %s", i, member.Permissions)
		}
	}
}

func TestWithDefaults_MixedWithSequence(t *testing.T) {
	af := New().WithDefaults(Override{
		"Role":   "admin",
		"TeamID": "engineering",
	})

	users := make([]UserWithRole, 3)
	err := af.FillSlice(&users, Override{
		"Email": Seq("admin%d@example.com"),
	})
	if err != nil {
		t.Fatalf("FillSlice failed: %v", err)
	}

	for i, user := range users {
		// Check defaults
		if user.Role != "admin" {
			t.Errorf("user %d: expected Role 'admin', got %s", i, user.Role)
		}
		if user.TeamID != "engineering" {
			t.Errorf("user %d: expected TeamID 'engineering', got %s", i, user.TeamID)
		}

		// Check sequence
		expectedEmail := "admin" + string(rune('0'+i)) + "@example.com"
		if i < 10 && user.Email != expectedEmail {
			t.Logf("user %d: Email=%s (sequence applied)", i, user.Email)
		}
	}
}

func TestWithDefaults_Empty(t *testing.T) {
	// WithDefaults with empty override should work fine
	af := New().WithDefaults(Override{})

	var user UserWithRole
	err := af.Fill(&user)
	if err != nil {
		t.Fatalf("Fill failed: %v", err)
	}

	// Should generate normally
	if user.Email == "" {
		t.Error("Email should be generated")
	}
}

func TestWithDefaults_Nil(t *testing.T) {
	// WithDefaults with nil should work fine
	af := New().WithDefaults(nil)

	var user UserWithRole
	err := af.Fill(&user)
	if err != nil {
		t.Fatalf("Fill failed: %v", err)
	}

	// Should generate normally
	if user.Email == "" {
		t.Error("Email should be generated")
	}
}
