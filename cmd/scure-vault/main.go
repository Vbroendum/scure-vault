package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	// This imports the cryptography package you just saved
	"github.com/Vbroendum/scure-vault/internal/crypto"
)

// Helper function to update .gitignore safely
func updateGitignore() error {
	ignoreFile := ".gitignore"
	// Rules we want to ensure exist
	rules := []string{".env", ".env.keys", "scure-vault.exe", "scure-vault"}

	// Read existing .gitignore (or create empty if missing)
	content, err := os.ReadFile(ignoreFile)
	existing := ""
	if err == nil {
		existing = string(content)
	}

	f, err := os.OpenFile(ignoreFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Append rules only if they are missing
	for _, rule := range rules {
		if !strings.Contains(existing, rule) {
			if _, err := f.WriteString(rule + "\n"); err != nil {
				return err
			}
			fmt.Printf("📝 Added %s to .gitignore\n", rule)
		}
	}
	return nil
}

// hello

func main() {
	var rootCmd = &cobra.Command{
		Use:   "scure-vault",
		Short: "A custom secure secrets manager CLI",
	}

	var pushCmd = &cobra.Command{
		Use:   "push",
		Short: "Encrypt the local .env into a .env.vault file",
		Run: func(cmd *cobra.Command, args []string) {
			key := crypto.GetMasterKey()
			err := crypto.EncryptFile(".env", ".env.vault", key)
			if err != nil {
				fmt.Printf("❌ Push failed: %v\n", err)
				return
			}
			fmt.Println("🔒 Successfully encrypted .env to .env.vault")
		},
	}

	var pullCmd = &cobra.Command{
		Use:   "pull",
		Short: "Decrypt the .env.vault file back into a .env file",
		Run: func(cmd *cobra.Command, args []string) {
			key := crypto.GetMasterKey()
			err := crypto.DecryptFile(".env.vault", ".env", key)
			if err != nil {
				fmt.Printf("❌ Pull failed: %v\n", err)
				return
			}
			fmt.Println("🔓 Successfully decrypted .env.vault to .env")
		},
	}

	var editCmd = &cobra.Command{
		Use:   "edit",
		Short: "Securely edit the encrypted vault and update .env",
		Run: func(cmd *cobra.Command, args []string) {
			key := crypto.GetMasterKey()

			// 1. Create temp file
			tmpFile, err := os.CreateTemp("", "scure-vault-*.env")
			if err != nil {
				fmt.Printf("❌ Failed to create temp file: %v\n", err)
				return
			}
			tmpName := tmpFile.Name()
			defer os.Remove(tmpName) // Clean up temp file at the end
			tmpFile.Close()

			// 2. Decrypt current vault to temp file
			// If vault exists, load it. If not, we start empty.
			if _, err := os.Stat(".env.vault"); err == nil {
				err := crypto.DecryptFile(".env.vault", tmpName, key)
				if err != nil {
					fmt.Printf("❌ Failed to decrypt vault: %v\n", err)
					return
				}
			}

			// 3. Determine Editor
			editor := os.Getenv("EDITOR")
			if editor == "" {
				if os.PathSeparator == '\\' {
					editor = "notepad"
				} else {
					editor = "nano"
				}
			}

			fmt.Printf("📝 Opening vault in %s... (Waiting for you to save and close)\n", editor)

			// 4. Run the editor
			parts := strings.Fields(editor)
			head := parts[0]
			editorArgs := parts[1:]
			editorArgs = append(editorArgs, tmpName)

			runEditor := exec.Command(head, editorArgs...)
			runEditor.Stdin = os.Stdin
			runEditor.Stdout = os.Stdout
			runEditor.Stderr = os.Stderr

			err = runEditor.Run()
			if err != nil {
				fmt.Printf("❌ Editor failed: %v\n", err)
				return
			}

			// 5. Encrypt the temp file back to .env.vault
			err = crypto.EncryptFile(tmpName, ".env.vault", key)
			if err != nil {
				fmt.Printf("❌ Failed to save changes to vault: %v\n", err)
				return
			}

			// 6. UPDATE THE LOCAL .ENV FILE
			// We read the temp file (which has your edits) and write it to .env
			newContent, err := os.ReadFile(tmpName)
			if err != nil {
				fmt.Printf("❌ Failed to read temp file: %v\n", err)
				return
			}

			err = os.WriteFile(".env", newContent, 0644)
			if err != nil {
				fmt.Printf("❌ Failed to update local .env: %v\n", err)
				return
			}

			fmt.Println("🔒 Vault updated (.env.vault)")
			fmt.Println("✅ Local .env updated")
		},
	}

	var viewCmd = &cobra.Command{
		Use:   "view",
		Short: "Safely view the .env file with values cloaked",
		Run: func(cmd *cobra.Command, args []string) {
			file, err := os.Open(".env")
			if err != nil {
				fmt.Println("❌ No .env file found in this directory.")
				return
			}
			defer file.Close()

			fmt.Println("👀 Viewing .env (Auto-cloaked):")
			fmt.Println("--------------------------------")

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				trimmed := strings.TrimSpace(line)

				if trimmed == "" || strings.HasPrefix(trimmed, "#") {
					fmt.Println(line)
					continue
				}

				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					fmt.Printf("\033[36m%s\033[0m=████████████████\n", parts[0])
				} else {
					fmt.Println(line)
				}
			}
		},
	}

	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate a secure master key and save it to .env.keys",
		Run: func(cmd *cobra.Command, args []string) {
			err := crypto.GenerateKey()
			if err != nil {
				fmt.Println("❌ Failed to generate key:", err)
				return
			}
			fmt.Println("✅ Generated new master key and saved to .env.keys")
			fmt.Println("⚠️  CRITICAL: Add .env.keys to your .gitignore immediately!")
		},
	}

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize the project: generates key, updates gitignore, and encrypts .env",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🚀 Initializing scure-vault...")

			// 1. Generate the Key
			if _, err := os.Stat(".env.keys"); err == nil {
				fmt.Println("⚠️  .env.keys already exists. Skipping key generation to prevent overwriting.")
			} else {
				err := crypto.GenerateKey()
				if err != nil {
					fmt.Printf("❌ Failed to generate key: %v\n", err)
					return
				}
				fmt.Println("✅ Generated secure master key in .env.keys")
			}

			// 2. Update .gitignore
			if err := updateGitignore(); err != nil {
				fmt.Printf("❌ Failed to update .gitignore: %v\n", err)
			}

			// 3. Encrypt the .env file (Push)
			if _, err := os.Stat(".env"); os.IsNotExist(err) {
				fmt.Println("⚠️  No .env file found. Creating an empty one for you.")
				os.WriteFile(".env", []byte("# Add your secrets here\n"), 0644)
			}

			key := crypto.GetMasterKey()
			err := crypto.EncryptFile(".env", ".env.vault", key)
			if err != nil {
				fmt.Printf("❌ Failed to encrypt vault: %v\n", err)
				return
			}
			fmt.Println("🔒 Encrypted .env to .env.vault")

			// 4. Show the Key to the User
			keyData, _ := os.ReadFile(".env.keys")
			keyString := strings.TrimSpace(string(keyData))

			fmt.Println("\n🎉 Initialization Complete!")
			fmt.Println("---------------------------------------------------")
			fmt.Println("🔑 HERE IS YOUR MASTER KEY (SAVE THIS SECURELY):")
			fmt.Printf("\n   \033[32m%s\033[0m\n\n", strings.TrimPrefix(keyString, "VAULT_KEY="))
			fmt.Println("Share this key with your team via 1Password/Signal.")
			fmt.Println("They will put it in their own .env.keys file.")
			fmt.Println("---------------------------------------------------")
		},
	}

	// Wire up ALL commands BEFORE executing
	rootCmd.AddCommand(pushCmd, pullCmd, viewCmd, generateCmd, initCmd, editCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
