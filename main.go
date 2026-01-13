package main

import (
	"flag"
	"github.com/Helland369/gitstats/git_stats"
)

func main() {
	var folder string
	var email string
	flag.StringVar(&folder, "add", "", "add a new folder to scan Git repositories")
	flag.StringVar(&email, "email", "thomas_helland@pm.me", "the email to scan")
	flag.Parse()

	if folder != "" {
		git_stats.Scan(folder)
		return
	}
	git_stats.Stats(email)
}
