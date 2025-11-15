package main

import (
	"testing"

	"github.com/m1a9s9a4/autofill"
)

// TestFillEntStruct demonstrates filling an Ent-generated struct
func TestFillEntStruct(t *testing.T) {
	// 型定義不要！Entの構造体を直接使う
	var user User
	err := autofill.Fill(&user, autofill.Override{
		"Email": "test@example.com",
		"Age":   25,
	})

	if err != nil {
		t.Fatalf("Fill failed: %v", err)
	}

	// 値の検証
	if user.Email != "test@example.com" {
		t.Errorf("expected Email test@example.com, got %s", user.Email)
	}
	if user.Age != 25 {
		t.Errorf("expected Age 25, got %d", user.Age)
	}
	if user.Name == "" {
		t.Error("Name should be generated")
	}
}

// TestCreateMultipleUsers demonstrates creating multiple test users
func TestCreateMultipleUsers(t *testing.T) {
	users := make([]User, 5)
	err := autofill.FillSlice(&users, autofill.Override{
		"Email":    autofill.Seq("user%d@example.com"),
		"TenantID": int64(1000),
		"Active":   true,
	})

	if err != nil {
		t.Fatalf("FillSlice failed: %v", err)
	}

	// 全ユーザーが同じTenantID
	for i, user := range users {
		if user.TenantID != 1000 {
			t.Errorf("user %d: expected TenantID 1000, got %d", i, user.TenantID)
		}
		if !user.Active {
			t.Errorf("user %d: expected Active true", i)
		}
	}
}

// TestRoleBasedFillers demonstrates creating role-specific fillers
func TestRoleBasedFillers(t *testing.T) {
	// Adminフィラー
	adminFiller := autofill.New().WithDefaults(autofill.Override{
		"Role":     "admin",
		"Active":   true,
		"TenantID": int64(1000),
	})

	// Memberフィラー
	memberFiller := autofill.New().WithDefaults(autofill.Override{
		"Role":     "member",
		"Active":   true,
		"TenantID": int64(1000),
	})

	// Admins作成
	admins := make([]User, 2)
	if err := adminFiller.FillSlice(&admins); err != nil {
		t.Fatalf("adminFiller failed: %v", err)
	}

	// Members作成
	members := make([]User, 3)
	if err := memberFiller.FillSlice(&members); err != nil {
		t.Fatalf("memberFiller failed: %v", err)
	}

	// 検証
	for i, admin := range admins {
		if admin.Role != "admin" {
			t.Errorf("admin %d: expected Role admin, got %s", i, admin.Role)
		}
	}

	for i, member := range members {
		if member.Role != "member" {
			t.Errorf("member %d: expected Role member, got %s", i, member.Role)
		}
	}
}

// TestRelatedEntities demonstrates creating related entities
func TestRelatedEntities(t *testing.T) {
	// Author作成
	var author User
	err := autofill.Fill(&author, autofill.Override{
		"ID":   int64(100),
		"Name": "Test Author",
	})
	if err != nil {
		t.Fatalf("Fill author failed: %v", err)
	}

	// Authorの投稿を複数作成
	posts := make([]Post, 3)
	err = autofill.FillSlice(&posts, autofill.Override{
		"AuthorID": author.ID, // 全てのPostが同じAuthor
		"Status":   "published",
	})
	if err != nil {
		t.Fatalf("FillSlice posts failed: %v", err)
	}

	// 全PostのAuthorIDが同じか確認
	for i, post := range posts {
		if post.AuthorID != author.ID {
			t.Errorf("post %d: expected AuthorID %d, got %d",
				i, author.ID, post.AuthorID)
		}
		if post.Status != "published" {
			t.Errorf("post %d: expected Status published, got %s",
				i, post.Status)
		}
	}
}

// BenchmarkFillEntStruct benchmarks filling an Ent struct
func BenchmarkFillEntStruct(b *testing.B) {
	af := autofill.New()
	var user User

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		af.Fill(&user)
	}
}

// Helper function example: Create test user with defaults
func CreateTestUser(t *testing.T, overrides ...autofill.Override) *User {
	t.Helper()

	var user User
	err := autofill.Fill(&user, overrides...)
	if err != nil {
		t.Fatalf("CreateTestUser failed: %v", err)
	}
	return &user
}

// TestHelperFunction demonstrates using helper functions
func TestHelperFunction(t *testing.T) {
	// シンプルに使える
	user1 := CreateTestUser(t)
	if user1.Name == "" {
		t.Error("user1 Name should not be empty")
	}

	// オーバーライドもできる
	user2 := CreateTestUser(t, autofill.Override{
		"Email": "custom@example.com",
		"Age":   40,
	})
	if user2.Email != "custom@example.com" {
		t.Errorf("expected custom@example.com, got %s", user2.Email)
	}
}
