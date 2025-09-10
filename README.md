<h1 align="center">STILLHERE - Dead Man's Switch</h1>

<p align="center">
  <img src="stillhere.png" alt="STILLHERE logo" width="300">
</p>

---

**stillhere** is a redundant, automated Dead Man's Switch powered by GitHub Actions. It periodically sends reminders and can trigger a final payload delivery (like an encrypted file) if you fail to check in after a configured number of days.

The system supports multiple GitHub accounts and repositories for added redundancy.

---

## ✨ Features

* ✅ Daily check-in via CLI tool (Linux, macOS, Windows) or CURL
* ✅ Redundant setup across multiple GitHub repositories
* ✅ Email reminders after missing check-ins
* ✅ Configurable payload release after extended absence

---

## ⚙️ How It Works

1. **Check-in**
   - Run the `checkin` tool daily to update the `last_checkin.txt` file across all configured repositories
   - Multiple check-ins on the same day are automatically ignored

2. **Daily Monitoring**
   - A scheduled GitHub Action runs daily at a specified time
   - Calculates days since last check-in

3. **Reminders**
   - If you miss check-ins for X days, the system sends email reminders

4. **Payload Delivery**
   - After Y days without check-in, the system sends a final email to designated recipients containing:
     - Last check-in date
     - Payload download instructions
     - Decryption key (or partial key for split-secret setups)

## 📂 Repository Structure

```
.
├── .github/
│   └── workflows/
│       ├── monitor.yml          # Monitors check-ins and triggers alerts
│       ├── checkin.yml          # Handles check-in operations
│       └── release.yml          # Manages release creation
├── cmd/
│   └── checkin/
│       └── main.go              # Check-in tool source code
├── last_checkin.txt             # Tracks latest check-in timestamp
├── secret.txt.enc               # Example encrypted payload
├── .env.example                 # Environment variables template
└── README.md
```

## 🚀 Setup Guide

### 1. Repository Setup
- Use this template to create one or more **private repositories**
- For maximum redundancy, set up across multiple GitHub accounts

### 2. Configure Secrets & Variables
Navigate to: `Settings → Secrets and variables → Actions` in each repository

**Required Secrets:**
```
SMTP_USER=your_email@gmail.com
SMTP_PASS=your_app_password
DECRYPTION_KEY=your_encryption_key
```

**Required Variables:**
```
USER_EMAIL=your_email@example.com
RECIPIENT_EMAIL=recipient@example.com
PAYLOAD_URL=https://example.com/payload.enc
FILE_NAME=payload.enc
REMINDER_DAY=3
RELEASE_DAY=7
```

> [!TIP]
> Choose either:
> - Attach encrypted files directly to emails, or
> - Host externally and provide download links via `PAYLOAD_URL`

### 3. Customize Email Content
Edit `.github/workflows/dms.yml` to modify email templates and messaging.

The system supports **split-key encryption** for enhanced security:
```bash
openssl aes-256-cbc -d -pbkdf2 -in <FILE> -out secret.txt \
  -pass pass:<PERSONAL_KNOWLEDGE>_${{ secrets.DECRYPTION_KEY }}_<OTHER_INFO>
```

### 4. Encrypt Your Payload
```bash
openssl aes-256-cbc -pbkdf2 -in secret.txt -out secret.txt.enc \
  -pass pass:your_encryption_key_here
```

### 5. Build Check-in Tool
Pre-built binaries are available through GitHub Actions:
1. Navigate to **Actions** tab
2. Run **Build Check-in Binary** workflow
3. Download binaries from the Releases page

Or build locally:
```bash
cd cmd/checkin
go build -o checkin
```

### 6. Configure Environment
Create a `.env` file:

**Single repository:**
```ini
GITHUB_REPOS=username/repository-name
GITHUB_TOKENS=github_personal_access_token
```

**Multiple repositories (comma-separated):**
```ini
GITHUB_REPOS=user1/repo1,user2/repo2,user3/repo3
GITHUB_TOKENS=token1,token2,token3
```

> [!WARNING]
> Number of repositories must match number of tokens

### 7. Daily Check-in
Run once daily:
```bash
./checkin
```
This updates all configured repositories with the current timestamp.

## 🔒 Security Considerations

- Maintain repositories as **private**
- Use [app passwords](https://support.google.com/accounts/answer/185833) for SMTP authentication
- Implement split-secret encryption for additional protection
- Regularly rotate credentials and tokens
- Disguise encrypted payloads as innocuous files
- Consider using hardware security keys for critical accounts

