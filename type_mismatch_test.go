package autofill

import (
	"testing"
)

type UserWithUUID struct {
	ID    string `autofill:"uuid"` // UUID (string)
	Name  string
	Age   int
	Email string
}

func TestTypeMismatch_IntToString(t *testing.T) {
	// IDはstring（UUID）だが、int64を指定
	var user UserWithUUID
	err := Fill(&user, Override{
		"ID": int64(12345), // 型不一致: int64 -> string
	})

	// エラーが発生するべき
	if err == nil {
		t.Fatalf("expected error for int64 -> string conversion, got nil. ID value: %q", user.ID)
	}

	t.Logf("✅ Correctly rejected int64 -> string: %v", err)
}

func TestTypeMismatch_StringToInt(t *testing.T) {
	// Ageはintだが、stringを指定
	var user UserWithUUID
	err := Fill(&user, Override{
		"Age": "not a number", // 型不一致: string -> int
	})

	// エラーが発生するべき
	if err == nil {
		t.Fatal("expected error for type mismatch, got nil")
	}
	t.Logf("✅ Correctly rejected string -> int: %v", err)
}

func TestTypeMatch_CorrectTypes(t *testing.T) {
	// 正しい型を指定すればエラーなし
	var user UserWithUUID
	err := Fill(&user, Override{
		"ID":    "custom-uuid-12345", // string -> string: OK
		"Age":   30,                   // int -> int: OK
		"Email": "test@example.com",   // string -> string: OK
	})

	if err != nil {
		t.Fatalf("unexpected error with correct types: %v", err)
	}

	if user.ID != "custom-uuid-12345" {
		t.Errorf("expected ID 'custom-uuid-12345', got %s", user.ID)
	}
	if user.Age != 30 {
		t.Errorf("expected Age 30, got %d", user.Age)
	}
}

func TestTypeConversion_CompatibleTypes(t *testing.T) {
	type Numbers struct {
		IntField   int
		Int64Field int64
		FloatField float64
	}

	var nums Numbers

	// int -> int64は変換可能
	err := Fill(&nums, Override{
		"Int64Field": int(100), // int -> int64: 変換可能
	})

	if err != nil {
		t.Logf("Conversion result: %v", err)
	}

	// int64 -> intは変換できるか試す
	err = Fill(&nums, Override{
		"IntField": int64(200), // int64 -> int: 変換可能
	})

	if err != nil {
		t.Logf("int64 -> int conversion: %v", err)
	}
}

func TestWithDefaults_TypeMismatch(t *testing.T) {
	// WithDefaultsで型不一致の値を設定
	af := New().WithDefaults(Override{
		"ID": int64(12345), // 型不一致
	})

	var user UserWithUUID
	err := af.Fill(&user)

	// エラーが発生するべき
	if err == nil {
		t.Fatal("expected error for type mismatch in WithDefaults, got nil")
	}
	t.Logf("Expected error with WithDefaults: %v", err)
}

func TestSlice_TypeMismatch(t *testing.T) {
	users := make([]UserWithUUID, 3)
	err := FillSlice(&users, Override{
		"ID": int64(12345), // 全要素に型不一致の値
	})

	// エラーが発生するべき
	if err == nil {
		t.Fatal("expected error for type mismatch in FillSlice, got nil")
	}
	t.Logf("Expected error in FillSlice: %v", err)
}
