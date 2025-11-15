# Ent Integration Example

This example demonstrates how to use autofill with Ent-generated structs.

## Setup

```bash
# Initialize Ent
go get entgo.io/ent/cmd/ent
ent init User

# Generate code
go generate ./ent
```

## Usage

```go
import (
    "your-project/ent"
    "github.com/m1a9s9a4/autofill"
)

// ✅ Method 1: Fill existing Ent struct
var user ent.User
autofill.Fill(&user, autofill.Override{
    "Email": "test@example.com",
})

// ✅ Method 2: Create helper function
func CreateTestUser(t *testing.T) *ent.User {
    var user ent.User
    autofill.Fill(&user)
    return &user
}

// ✅ Method 3: With DB integration
func CreateTestUserInDB(t *testing.T, client *ent.Client) *ent.User {
    var user ent.User
    autofill.Fill(&user)

    return client.User.Create().
        SetName(user.Name).
        SetEmail(user.Email).
        SetAge(user.Age).
        Save(context.Background())
}
```

## No Need to Define Types!

You don't need to define your own structs - just use what Ent generates:

```go
// ❌ Don't do this
type User struct {
    Name  string
    Email string
}

// ✅ Do this instead
import "your-project/ent"
var user ent.User  // Use Ent's generated struct
autofill.Fill(&user)
```
