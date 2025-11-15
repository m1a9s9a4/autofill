package main

import (
	"fmt"

	"github.com/m1a9s9a4/autofill"
)

type User struct {
	ID          int64  `autofill:"seq"`
	Email       string `autofill:"email"`
	Name        string
	Role        string
	Permissions string
	WorkspaceID int64
}

func main() {
	fmt.Println("=== WithDefaults: Admin vs Member Example ===")

	// Admin用のフィラーを作成
	adminFiller := autofill.New().
		WithSeed(12345).
		WithDefaults(autofill.Override{
			"Role":        "admin",
			"Permissions": "all",
			"WorkspaceID": int64(1000),
		})

	// メンバー用のフィラーを作成
	memberFiller := autofill.New().
		WithSeed(12345).
		WithDefaults(autofill.Override{
			"Role":        "member",
			"Permissions": "read",
			"WorkspaceID": int64(1000),
		})

	// Adminユーザーを3人作成
	admins := make([]User, 3)
	adminFiller.FillSlice(&admins, autofill.Override{
		"Email": autofill.Seq("admin%d@example.com"),
	})

	// 一般メンバーを5人作成
	members := make([]User, 5)
	memberFiller.FillSlice(&members, autofill.Override{
		"Email": autofill.Seq("member%d@example.com"),
	})

	// 結果表示
	fmt.Println("📋 Admins:")
	for i, admin := range admins {
		fmt.Printf("  %d: ID=%d Email=%s Role=%s Permissions=%s WorkspaceID=%d\n",
			i+1, admin.ID, admin.Email, admin.Role, admin.Permissions, admin.WorkspaceID)
	}

	fmt.Println("\n👥 Members:")
	for i, member := range members {
		fmt.Printf("  %d: ID=%d Email=%s Role=%s Permissions=%s WorkspaceID=%d\n",
			i+1, member.ID, member.Email, member.Role, member.Permissions, member.WorkspaceID)
	}

	fmt.Println("\n=== Override Defaults Example ===")

	// デフォルトはadminだが、一部だけsuperadminにする
	users := make([]User, 3)
	adminFiller.Fill(&users[0]) // admin (デフォルト)
	adminFiller.Fill(&users[1], autofill.Override{
		"Role":        "superadmin", // オーバーライド
		"Permissions": "unlimited",
	})
	adminFiller.Fill(&users[2]) // admin (デフォルト)

	fmt.Println("👑 Mixed Roles:")
	for i, user := range users {
		fmt.Printf("  %d: Role=%s Permissions=%s\n",
			i+1, user.Role, user.Permissions)
	}

	fmt.Println("\n=== Different Workspaces Example ===")

	// ワークスペースAのユーザー
	workspaceAFiller := autofill.New().WithDefaults(autofill.Override{
		"WorkspaceID": int64(100),
		"Role":        "member",
	})

	// ワークスペースBのユーザー
	workspaceBFiller := autofill.New().WithDefaults(autofill.Override{
		"WorkspaceID": int64(200),
		"Role":        "member",
	})

	workspaceAUsers := make([]User, 2)
	workspaceBUsers := make([]User, 2)

	workspaceAFiller.FillSlice(&workspaceAUsers)
	workspaceBFiller.FillSlice(&workspaceBUsers)

	fmt.Println("🏢 Workspace A (ID=100):")
	for i, user := range workspaceAUsers {
		fmt.Printf("  %d: ID=%d WorkspaceID=%d\n", i+1, user.ID, user.WorkspaceID)
	}

	fmt.Println("\n🏢 Workspace B (ID=200):")
	for i, user := range workspaceBUsers {
		fmt.Printf("  %d: ID=%d WorkspaceID=%d\n", i+1, user.ID, user.WorkspaceID)
	}
}
