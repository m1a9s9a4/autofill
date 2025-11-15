package main

import (
	"fmt"

	"github.com/m1a9s9a4/autofill"
)

// Entが生成する構造体のサンプル
// 実際のプロジェクトでは "your-project/ent" からimportする
type User struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Age      int    `json:"age,omitempty"`
	Active   bool   `json:"active,omitempty"`
	Role     string `json:"role,omitempty"`
	TenantID int64  `json:"tenant_id,omitempty"`
}

type Post struct {
	ID       int64  `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Content  string `json:"content,omitempty"`
	AuthorID int64  `json:"author_id,omitempty"`
	Status   string `json:"status,omitempty"`
}

func main() {
	fmt.Println("=== Ent Integration Examples ===")

	// Example 1: 型定義不要 - Entの構造体を直接使う
	fmt.Println("1. Using Ent-generated struct directly:")
	var user User
	autofill.Fill(&user, autofill.Override{
		"Email": "john.doe@example.com",
		"Age":   30,
	})
	fmt.Printf("   User: %+v\n", user)

	// Example 2: テストデータを複数作成
	fmt.Println("\n2. Creating multiple test users:")
	users := make([]User, 3)
	autofill.FillSlice(&users, autofill.Override{
		"Email":    autofill.Seq("user%d@example.com"),
		"TenantID": int64(1000), // 全ユーザー同じテナント
	})
	for i, u := range users {
		fmt.Printf("   User %d: Name=%s Email=%s TenantID=%d\n", i+1, u.Name, u.Email, u.TenantID)
	}

	// Example 3: ロール別のフィラーを作成
	fmt.Println("\n3. Role-specific fillers (WithDefaults):")

	adminFiller := autofill.New().WithDefaults(autofill.Override{
		"Role":     "admin",
		"Active":   true,
		"TenantID": int64(1000),
	})

	memberFiller := autofill.New().WithDefaults(autofill.Override{
		"Role":     "member",
		"Active":   true,
		"TenantID": int64(1000),
	})

	admins := make([]User, 2)
	adminFiller.FillSlice(&admins)

	members := make([]User, 3)
	memberFiller.FillSlice(&members)

	fmt.Println("   Admins:")
	for i, admin := range admins {
		fmt.Printf("     %d: Name=%s Role=%s\n", i+1, admin.Name, admin.Role)
	}

	fmt.Println("   Members:")
	for i, member := range members {
		fmt.Printf("     %d: Name=%s Role=%s\n", i+1, member.Name, member.Role)
	}

	// Example 4: 関連データの作成（PostとUser）
	fmt.Println("\n4. Creating related entities:")

	// ユーザー作成
	var author User
	autofill.Fill(&author, autofill.Override{
		"ID":    int64(100),
		"Name":  "Author Name",
		"Email": "author@example.com",
	})

	// そのユーザーの投稿を作成
	posts := make([]Post, 3)
	autofill.FillSlice(&posts, autofill.Override{
		"AuthorID": author.ID,   // 同じAuthorID
		"Status":   "published", // 全て公開済み
		"Title":    autofill.Seq("Post #%d"),
	})

	fmt.Printf("   Author: %s (ID=%d)\n", author.Name, author.ID)
	fmt.Println("   Posts:")
	for i, post := range posts {
		fmt.Printf("     %d: Title=%s AuthorID=%d Status=%s\n",
			i+1, post.Title, post.AuthorID, post.Status)
	}

	fmt.Println("\n=== Summary ===")
	fmt.Println("✅ No need to define types - use Ent's generated structs")
	fmt.Println("✅ Fill with default values using WithDefaults")
	fmt.Println("✅ Override specific fields as needed")
	fmt.Println("✅ Perfect for testing and seeding data")
}
