package main

import (
	"flag"
	
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
		github_stats.Get_github_contrib(userName)		
	}
	
	git_stats.Stats(email)
}
