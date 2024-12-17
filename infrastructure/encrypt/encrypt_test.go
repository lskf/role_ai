package encrypt

import "testing"

func TestGenerateKey(t *testing.T) {
	// TestGenerateKey is a test function.
	t.Run("TestGenerateKey", func(t *testing.T) {
		// TestGenerateKey is a test function.
		key, err := GenerateEncryptKey()
		if err != nil {
			t.Errorf("GenerateEncryptKey() error = %v", err)
			return
		}
		t.Logf("GenerateEncryptKey() = %v", key)
	})
}

func TestEncrypt(t *testing.T) {
	// TestEncrypt is a test function.
	t.Run("TestEncrypt", func(t *testing.T) {
		// TestEncrypt is a test function.
		encrypted, err := Encrypt("test", "GfDrRibIFEa0ErmZpjZ2rQvcSEHi1igartDmPi2GpP4=")
		if err != nil {
			t.Errorf("Encrypt() error = %v", err)
			return
		}
		t.Logf("Encrypt() = %v", encrypted)
	})
}

func TestDecrypt(t *testing.T) {
	// TestDecrypt is a test function.
	t.Run("TestDecrypt", func(t *testing.T) {
		// TestDecrypt is a test function.
		decrypted, err := Decrypt("7DlhrMNztBTgrkis7oik6bJI5YNzu6WajNpICSugtIw=", "GfDrRibIFEa0ErmZpjZ2rQvcSEHi1igartDmPi2GpP4=")
		if err != nil {
			t.Errorf("Decrypt() error = %v", err)
			return
		}
		t.Logf("Decrypt() = %v", decrypted)
	})
}
