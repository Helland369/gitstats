package github_stats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const url = "https://api.github.com/graphql"

type graphQLRequest struct {
	Query string `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

func Get_github_contrib(userName string) ([]byte, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found!")
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		panic("set GITHUB_TOKE env variable to github personal access token!")
	}

	query := `
           query($login: String!) {
             user(login: $login) {
               contributionsCollection {
                 contributionCalendar {
                   totalContributions
                   weeks {
                     firstDay
                     contributionDays {
                       date
                       weekday
                       contributionCount
                       color
                     }
                   }
                 }
               }
             }
           }`

	requBody, err := json.Marshal(graphQLRequest {
		Query: query,
		Variables: map[string]any{
			"login": userName,
		},
	})

	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewReader(requBody))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer " + token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github+json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	
	if res.StatusCode != 200 {
		return body, fmt.Errorf("HTTP %d: %s", res.StatusCode, string(body))
	}

	return body, nil
}
