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

func Get_github_contrib(userName string) {
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

	requBody, _ := json.Marshal(graphQLRequest {
		Query: query,
		Variables: map[string]any{
			"login": userName,
		},
	})

	req, _ := http.NewRequest("POST", url, bytes.NewReader(requBody))
	req.Header.Set("Authorization", "Bearer " + token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github+json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	body, err := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		fmt.Printf("HTTP %d\n%s\n", res.StatusCode, string(body))
		return
	}
	if err != nil {
		fmt.Println("Failed to read response body: ", err)
		return
	}

	fmt.Println(string(body))
}
