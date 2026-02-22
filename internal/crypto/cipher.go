package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

// GetMasterKey fetches the decryption key from the environment or .env.keys file
func GetMasterKey() []byte {
	keyHex := os.Getenv("VAULT_KEY")

	// If not in system env, check the .env.keys file
	if keyHex == "" {
		data, err := os.ReadFile(".env.keys")
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "VAULT_KEY=") {
					keyHex = strings.TrimSpace(strings.TrimPrefix(line, "VAULT_KEY="))
					break
				}
			}
		}
	}

	if keyHex == "" {
		fmt.Println("❌ Error: VAULT_KEY not found in system environment or .env.keys file.")
		fmt.Println("💡 Run '.\\secure-vault.exe generate' to create one.")
		os.Exit(1)
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil || len(key) != 32 {
		fmt.Println("❌ Error: VAULT_KEY must be a valid 64-character hex string (32 bytes).")
		os.Exit(1)
	}
	return key
}

// EncryptFile reads a plaintext file, encrypts it, and saves it.
func EncryptFile(inFile string, outFile string, key []byte) error {
	plaintext, err := os.ReadFile(inFile)
	if err != nil {
		return fmt.Errorf("could not read %s: %v", inFile, err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return os.WriteFile(outFile, ciphertext, 0644)
}

// DecryptFile reads the encrypted vault, extracts the nonce, decrypts, and saves it.
func DecryptFile(inFile string, outFile string, key []byte) error {
	ciphertext, err := os.ReadFile(inFile)
	if err != nil {
		return fmt.Errorf("could not read %s: %v", inFile, err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext is too short or corrupted")
	}

	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption failed (wrong key or corrupted file): %v", err)
	}

	return os.WriteFile(outFile, plaintext, 0644)
}

// GenerateKey creates a 32-byte master key and saves it to .env.keys
func GenerateKey() error {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return err
	}
	keyHex := hex.EncodeToString(key)
	fileContent := fmt.Sprintf("VAULT_KEY=%s\n", keyHex)
	return os.WriteFile(".env.keys", []byte(fileContent), 0600)
}