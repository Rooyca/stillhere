# stillhere-dms · Dead Man’s Switch

**stillhere-dms** is a **redundant, automated Dead Man’s Switch** powered by **GitHub Actions**.
It automatically sends reminders and can trigger a final payload (such as an encrypted file) if you don’t check in for a set number of days.

The system supports **multiple GitHub accounts/repos for redundancy**:

* If one repo is suspended, removed, or unavailable, the others continue functioning.

---

## ✨ Features

* ✅ **Daily check-in system** via CLI tool
* ✅ **Redundant setup** across one or more GitHub repos
* ✅ **Reminder emails** if no check-in for **3–6 days**
* ✅ **Payload release** after **7+ days** without check-in
* ✅ **Optional encrypted payload file** with flexible key handling
* ✅ **Cross-platform Go-based CLI tool** (Linux, macOS, Windows)

---

## ⚙️ How It Works

1. **Check-in**

   * Each day, you run the `checkin` binary.
   * This updates the `last_checkin.txt` file in all configured repos.
   * Multiple check-ins on the same day are ignored.

2. **Daily Monitoring**

   * A scheduled GitHub Action runs daily at a fixed time.
   * It calculates how many days have passed since your last check-in.

3. **Reminders**

   * If you haven’t checked in for **3–6 days**, you receive an email reminder.

4. **Payload Trigger**

   * If you haven’t checked in for **7+ days**, an email is sent to your designated recipient.
   * The email includes:

     * The date of the last check-in
     * Instructions on downloading and decrypting the payload
     * A decryption key (or part of it, if using split-secrets)

---

## 📂 Repo Structure

```
.
├── .github/
│   └── workflows/
│       ├── monitor.yml          # Workflow to send reminders and release payload
│       ├── checkin.yml          # Workflow to checkin
│       └── release.yml          # Workflow to create new release
├── cmd/
│   └── checkin/
│       └── main.go              # Go code for check-in tool
├── last_checkin.txt             # Updated daily by check-ins
├── secret.txt.enc               # Example encrypted payload file
├── .env.example                 # Example enviroment variables for checkin
└── README.md
```

---

## 🚀 Setup Guide

### 1. Create GitHub Repos

* Create one or more **private repositories** named `stillhere-dms` (or similar).
* For redundancy, you can set up multiple accounts and repos.

---

### 2. Add the Workflows

* Copy `checkin.yml` and `monitor.yml` into `.github/workflows/` in each repo.

---

### 3. Configure Secrets & Variables

> [!WARNING]
> All secrets and variables must be set in each repo.

Go to:
`Settings → Secrets and variables → Actions`

**Secrets**

```
SMTP_USER=example@gmail.com
SMTP_PASS=your-app-password
DECRYPTION_KEY=someStrongKeyHere
```

**Variables**

```
USER_EMAIL=youremail@example.com
RECIPIENT_EMAIL=recipient@example.com
PAYLOAD_URL=https://domain.com/file.enc
FILE_NAME=filename.enc
REMINDER_DAY=3
RELEASE_DAY=7
```

> [!TIP]
> You can either:
> * Attach the encrypted file directly to the email, or
> * Host it externally and use `PAYLOAD_URL` to share the download link.

---

### 4. Customize Email Content

Edit `.github/workflows/dms.yml` to adjust the email body.

The project supports **split key handling**. Example:

```bash
openssl aes-256-cbc -d -pbkdf2 -in <FILENAME> -out secret.txt \
  -pass pass:<YOUR_BIRTHDAY_MMDDYYYY>_${{ secrets.DECRYPTION_KEY }}_<MY_ID_NUMBER>
```

This forces the recipient to combine personal knowledge with the repo secret to decrypt.

---

### 5. Encrypt Your Payload

```bash
# Encrypt a file
openssl aes-256-cbc -pbkdf2 -in secret.txt -out secret.txt.enc -pass pass:someStrongKeyHere
```

---

### 6. Build the Check-in Tool

The repo includes a workflow that compiles binaries for Linux, macOS, and Windows.

1. Go to your repo → **Actions** tab.
2. Run **Build Check-in Binary**.
3. Download the artifacts.

Build locally:

```bash
cd cmd/checkin
go build -o checkin
```

Cross-compile:

```bash
GOOS=linux   GOARCH=amd64 go build -o checkin_linux
GOOS=darwin  GOARCH=arm64 go build -o checkin_macos
GOOS=windows GOARCH=amd64 go build -o checkin.exe
```

---

### 7. Configure `.env`

On your local machine, create a `.env` file:

```ini
# Single repo setup
GITHUB_REPOS=user1/stillhere-dms
GITHUB_TOKENS=ghp_token1

# Multiple repos for redundancy
GITHUB_REPOS=user1/repo1,user2/repo2,user3/repo3
GITHUB_TOKENS=ghp_token1,ghp_token2,ghp_token3
```

> [!WARNING]
> The number of repos must match the number of tokens.

---

### 8. Run Daily Check-in

Run once per day:

```bash
./checkin
```

This updates all configured repos with the current date.

---

## 🔒 Security Notes

* Keep repos **private**.
* Use **app passwords** for SMTP.
* Consider **split-secrets**: store part of the key in the repo and another part offline.
* Rotate credentials periodically.
* Disguise encrypted payloads as ordinary files.

---

## 🛡️ Redundancy

* The **Go binary checks in to all repos** listed in `.env`.
* If one repo is deleted, others continue working.
* Only one check-in per day is needed.

---

## ⚡ Example Email Flow

**Day 3–6: Reminder**

> Subject: Reminder – Dead Man’s Switch
> You haven’t checked in for 4 days. Please check in before day 7, or the payload will be released.

**Day 7+: Trigger**

> Subject: Dead Man’s Switch Triggered
> The switch has been triggered. Last check-in: 2025-09-02
>
> To unlock the file:
>
> 1. Download `secret.txt.enc` from the repo.
> 2. Run:
>
>    ```bash
>    openssl aes-256-cbc -d -pbkdf2 -in secret.txt.enc -out secret.txt -pass pass:YOUR_KEY_HERE
>    ```
>
> 🔑 Decryption key: `someStrongKeyHere`

---

## ✅ Summary

* **Daily run:** `./checkin`
* **Supports multiple repos/accounts**
* **Reminder emails:** after 3–6 days of no check-in
* **Payload release:** after 7+ days
* **Encrypted payload delivery:** stored in repo, disguised, decrypted with shared key
