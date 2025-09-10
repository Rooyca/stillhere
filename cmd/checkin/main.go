package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func githubCheckin(repo, token, today string) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/dispatches", repo)
	payload := map[string]interface{}{
		"event_type": "checkin",
		"client_payload": map[string]string{
			"date": today,
		},
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.everest-preview+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("❌ GitHub checkin failed for %s: %v\n", repo, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("✅ GitHub checkin to %s: %s\n", repo, resp.Status)
}

func main() {
	// Load .env file if present
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️  No .env file found, using system env")
	}

	today := time.Now().UTC().Format("2006-01-02")

	githubRepos := strings.Split(os.Getenv("GITHUB_REPOS"), ",")
	githubTokens := strings.Split(os.Getenv("GITHUB_TOKENS"), ",")

	if len(githubRepos) != len(githubTokens) {
		log.Fatal("❌ Number of repos and tokens must match")
	}

	for i, repo := range githubRepos {
		repo = strings.TrimSpace(repo)
		token := strings.TrimSpace(githubTokens[i])

		if repo == "" || token == "" {
			continue
		}

		githubCheckin(repo, token, today)
	}

	log.Println("🎉 All checkins completed for:", today)
}
