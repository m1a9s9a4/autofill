package main

import (
	"testing"

	"github.com/m1a9s9a4/autofill"
)

type TestUser struct {
	ID          int64  `autofill:"seq"`
	Email       string
	WorkspaceID int64
	TeamID      string
	Role        string
}

func TestFixedValues(t *testing.T) {
	users := make([]TestUser, 5)
	err := autofill.FillSlice(&users, autofill.Override{
		"Email":       autofill.Seq("user%d@example.com"),
		"WorkspaceID": int64(12345),
		"TeamID":      "engineering",
		"Role":        "developer",
	})

	if err != nil {
		t.Fatalf("FillSlice failed: %v", err)
	}

	// Check that all users have the same fixed values
	for i, user := range users {
		// Fixed values should be the same for all users
		if user.WorkspaceID != 12345 {
			t.Errorf("user %d: expected WorkspaceID 12345, got %d", i, user.WorkspaceID)
		}
		if user.TeamID != "engineering" {
			t.Errorf("user %d: expected TeamID 'engineering', got %s", i, user.TeamID)
		}
		if user.Role != "developer" {
			t.Errorf("user %d: expected Role 'developer', got %s", i, user.Role)
		}
	}

	// Verify all users have the same WorkspaceID
	firstWorkspaceID := users[0].WorkspaceID
	for i := 1; i < len(users); i++ {
		if users[i].WorkspaceID != firstWorkspaceID {
			t.Errorf("WorkspaceID should be same for all users, got %d and %d",
				firstWorkspaceID, users[i].WorkspaceID)
		}
	}
}

func TestMixedFixedAndSequential(t *testing.T) {
	type Task struct {
		ID        int64 `autofill:"seq"`
		ProjectID int64
		Priority  int
		Status    string
	}

	tasks := make([]Task, 3)
	err := autofill.FillSlice(&tasks, autofill.Override{
		"ProjectID": int64(999),
		"Priority":  autofill.SeqInt(1),
		"Status":    "pending",
	})

	if err != nil {
		t.Fatalf("FillSlice failed: %v", err)
	}

	for i, task := range tasks {
		// Fixed values
		if task.ProjectID != 999 {
			t.Errorf("task %d: expected ProjectID 999, got %d", i, task.ProjectID)
		}
		if task.Status != "pending" {
			t.Errorf("task %d: expected Status 'pending', got %s", i, task.Status)
		}

		// Sequential values
		expectedPriority := 1 + i
		if task.Priority != expectedPriority {
			t.Errorf("task %d: expected Priority %d, got %d", i, expectedPriority, task.Priority)
		}
	}
}
