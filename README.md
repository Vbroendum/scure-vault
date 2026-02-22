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

---

## 📦 Installation

`scure-vault` is a standalone binary. You do **not** need to install Go, Node, or Python to use it.

**Mac / Linux (Homebrew):**
```bash
brew tap yourusername/tap
brew install scure-vault