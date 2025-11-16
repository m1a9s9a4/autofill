package main

import (
	"fmt"

	"github.com/m1a9s9a4/autofill"
)

type User struct {
	ID    string `autofill:"uuid"` // UUID (string型)
	Name  string
	Age   int
	Score float64
}

func main() {
	fmt.Println("=== Type Safety Examples ===")

	// ✅ 正しい型: エラーなし
	fmt.Println("1. Correct types:")
	var user1 User
	err := autofill.Fill(&user1, autofill.Override{
		"ID":    "custom-uuid-12345", // string -> string: OK
		"Age":   30,                   // int -> int: OK
		"Score": 95.5,                 // float64 -> float64: OK
	})
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Printf("   ✅ Success: ID=%s, Age=%d, Score=%.1f\n", user1.ID, user1.Age, user1.Score)
	}

	// ❌ 型不一致: int64 -> string はエラー
	fmt.Println("\n2. Type mismatch (int64 -> string):")
	var user2 User
	err = autofill.Fill(&user2, autofill.Override{
		"ID": int64(12345), // ❌ int64 -> string: エラー
	})
	if err != nil {
		fmt.Printf("   ✅ Correctly rejected: %v\n", err)
	} else {
		fmt.Printf("   ❌ Should have failed but didn't\n")
	}

	// ❌ 型不一致: string -> int はエラー
	fmt.Println("\n3. Type mismatch (string -> int):")
	var user3 User
	err = autofill.Fill(&user3, autofill.Override{
		"Age": "not a number", // ❌ string -> int: エラー
	})
	if err != nil {
		fmt.Printf("   ✅ Correctly rejected: %v\n", err)
	} else {
		fmt.Printf("   ❌ Should have failed but didn't\n")
	}

	// ✅ 数値型間の変換: OK
	fmt.Println("\n4. Numeric type conversions (allowed):")
	var user4 User
	err = autofill.Fill(&user4, autofill.Override{
		"Age":   int64(25),   // int64 -> int: OK (数値型同士)
		"Score": float32(88), // float32 -> float64: OK (数値型同士)
	})
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Printf("   ✅ Success: Age=%d, Score=%.1f (numeric conversions work)\n", user4.Age, user4.Score)
	}

	// WithDefaultsでの型安全性
	fmt.Println("\n5. Type safety with WithDefaults:")
	filler := autofill.New().WithDefaults(autofill.Override{
		"ID": int64(99999), // ❌ デフォルト値が型不一致
	})

	var user5 User
	err = filler.Fill(&user5)
	if err != nil {
		fmt.Printf("   ✅ Correctly rejected in WithDefaults: %v\n", err)
	} else {
		fmt.Printf("   ❌ Should have failed but didn't\n")
	}

	fmt.Println("\n=== Summary ===")
	fmt.Println("✅ String <-> String: Allowed")
	fmt.Println("✅ Int <-> Int64: Allowed (numeric types)")
	fmt.Println("✅ Float32 <-> Float64: Allowed (numeric types)")
	fmt.Println("❌ Int/Int64 -> String: Rejected (would become Unicode)")
	fmt.Println("❌ String -> Int: Rejected (not parseable)")
	fmt.Println("\n💡 Always use the correct type in Override to avoid errors!")
}
