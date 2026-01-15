package main

import (
	"flag"
	"fmt"

	"github.com/Helland369/gitstats/git_stats"
	"github.com/Helland369/gitstats/github_stats"
)

func main() {
	var folder string
	var email string
	var userName string
	flag.StringVar(&folder, "add", "", "add a new folder to scan Git repositories")
	flag.StringVar(&email, "email", "", "the email to scan")
	flag.StringVar(&userName, "user", "", "github user name")
	flag.Parse()

	if folder != "" {
		git_stats.Scan(folder)
		return
	}

  if userName != "" {
		res, err := github_stats.Github_stats(userName)
		if err != nil {
			println(err)
			return
		}
		
		commits := github_stats.To_commit_map(res)

		fmt.Println("Github")
		git_stats.Print_commit_stats(commits)

		return
	}

	if email != "" {
		fmt.Println("Local Git")
		git_stats.Stats(email)
	}
}
