package main

import (
	"fmt"
	"time"

	"github.com/m1a9s9a4/autofill"
)

type User struct {
	ID        int64     `autofill:"seq"`
	Name      string
	Email     string    `autofill:"email"`
	Age       int       `autofill:"min=18,max=65"`
	Active    bool
	CreatedAt time.Time `autofill:"now"`
}

func main() {
	fmt.Println("=== Basic Fill ===")
	var user User
	af := autofill.New().WithSeed(12345)
	if err := af.Fill(&user); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n\n", user)

	// Fill multiple
	fmt.Println("=== Fill Multiple ===")
	users := make([]User, 3)
	if err := af.FillSlice(&users); err != nil {
		panic(err)
	}
	for i, u := range users {
		fmt.Printf("%d: ID=%d Email=%s Age=%d\n", i+1, u.ID, u.Email, u.Age)
	}
	fmt.Println()

	// With override
	fmt.Println("=== With Override ===")
	if err := af.Fill(&user, autofill.Override{
		"Name":  "John Doe",
		"Email": "john@example.com",
		"Age":   25,
	}); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n\n", user)

	// With sequence
	fmt.Println("=== With Sequence ===")
	users = make([]User, 3)
	if err := af.FillSlice(&users, autofill.Override{
		"Email": autofill.Seq("user%d@example.com"),
		"Age":   autofill.SeqInt(20),
	}); err != nil {
		panic(err)
	}
	for i, u := range users {
		fmt.Printf("%d: Email=%s Age=%d\n", i+1, u.Email, u.Age)
	}
}
