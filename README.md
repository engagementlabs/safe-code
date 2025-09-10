# Safe-Code: Telegram Ping Bot with Cryptographic Verification

A simple Telegram bot that responds with 'pong' when you send the '/ping' command.

## ⚖️ LICENSE NOTICE

**This code is provided for VERIFICATION and TRANSPARENCY purposes only.**

- ✅ You MAY: View, audit, suggest improvements, contribute
- ❌ You MAY NOT: Use, copy, modify, distribute, or incorporate into other projects

See [LICENSE](LICENSE) for full terms. This is NOT open source software.

## Setup

1. **Create a Telegram Bot**
   - Message [@BotFather](https://t.me/botfather) on Telegram
   - Use `/newbot` command and follow the instructions
   - Save the bot token you receive

2. **Configure Environment**
   - Copy `.env.example` to `.env`
   - Replace `your_bot_token_here` with your actual bot token
   
   Or set the environment variable directly:
   ```bash
   export TELEGRAM_BOT_TOKEN=your_bot_token_here
   ```

## Usage

### Option 1: Run with Go
```bash
go run main.go
```

### Option 2: Run with Docker
1. **Build the image**
   ```bash
   docker build -t telegram-ping-bot .
   ```

2. **Run the container**
   ```bash
   docker run -e TELEGRAM_BOT_TOKEN=your_bot_token_here telegram-ping-bot
   ```

## Test the Bot
- Find your bot on Telegram using the username you created
- Send `/ping` command
- The bot will respond with 'pong'

## Commands

- `/ping` - Bot responds with 'pong'

## Verification Endpoints

- `/health` - Health check (returns "OK")
- `/version` - Shows build information and GitHub commit
- `/attestation` - **Cryptographic proof of code integrity**

### Code Integrity Verification

To **cryptographically verify** that the running code matches this GitHub repository:

```bash
# Download and run the verification script
curl -O https://raw.githubusercontent.com/engagementlabs/safe-code/main/verify.py
python3 verify.py https://safe-code-p.engagelabs.org
```

### Example /attestation response:
```json
{
  "commit_sha": "abc123def456",
  "build_time": "2024-01-15T10:30:00Z",
  "source_hash": "d4e5f6a7b8c9...",
  "build_hash": "1a2b3c4d5e6f...",
  "github_url": "https://github.com/engagementlabs/safe-code/commit/abc123def456",
  "attestation": "9f8e7d6c5b4a...",
  "nitro_document": "eyJ0eXAiOiJKV1Q..."
}
```

### How Verification Works:

1. **Source Hash**: Downloads source from GitHub and calculates SHA256
2. **Build Hash**: Verifies the running binary matches the build
3. **Attestation**: Cryptographic signature proving integrity
4. **Nitro Enclave**: Hardware-level attestation (when available)

**This provides cryptographic proof that the running code is exactly what's in this repository.**

## Requirements

- Go 1.21+