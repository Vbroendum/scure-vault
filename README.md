# 🔐 scure-vault

A lightning-fast, zero-dependency CLI tool for syncing encrypted `.env` files across development teams and multiple machines. Built in Go using AES-256-GCM encryption.



## 🧠 How it Works

`scure-vault` splits your environment variables into two parts so you can safely use Git to sync your secrets without ever exposing them.



```text
┌─────────────────┐       Git Push       ┌─────────────────┐
│  .env.vault     │ ───────────────────▶ │     GitHub      │ (Public/Shared)
│ (Encrypted Safe)│                      │   Repository    │
└─────────────────┘                      └─────────────────┘
                                                  │
┌─────────────────┐    1Password / Slack          │ Git Pull
│  .env.keys      │ ───────────────────┐          ▼
│  (Master Key)   │                    │ ┌─────────────────┐
└─────────────────┘                    └▶│  Teammate's PC  │ (Decrypted locally)
  (Stays Local)                          └─────────────────┘

```

## 📦 Installation

`scure-vault` is a standalone CLI tool. Choose the installation method that best fits your operating system below.

### Option 1: Go Install (Recommended for Windows)
The most reliable way to install `scure-vault` across all platforms, especially on Windows, is using the Go toolchain. This compiles the binary directly on your machine, bypassing strict OS security blockers (like Windows Smart App Control) that flag unsigned internet downloads.

**Prerequisite:** You must have [Go installed](https://go.dev/doc/install).

1. Run the following command in your terminal:
   ```bash
   go install [github.com/Vbroendum/scure-vault/cmd/scure-vault@latest](https://github.com/Vbroendum/scure-vault/cmd/scure-vault@latest)
   ```

### Option 2: Homebrew (macOS / Linux)

If you are on a Mac or Linux machine with Homebrew, you can install the pre-compiled binary in seconds.

  Add the custom tap and install the tool:
    
  ```Bash

  brew tap Vbroendum/homebrew-tap
  brew install scure-vault
  ```
    
  Verify the installation:
  ```Bash
  scure-vault
  ```
