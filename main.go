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

	// work in progress
  if userName != "" {
		res, err := github_stats.Github_stats(userName)
		if err != nil {
			println(err)
		}
		cal := res.Data.User.ContributionsCollection.ContributionCalendar
		fmt.Println("Total: ", cal.TotalContributions)
		for _, w := range cal.Weeks {
			for _, d := range w.ContributionDays {
				if d.ContributionCount > 0 {
					fmt.Printf("%s (%d): %d\n", d.Date, d.Weekday, d.ContributionCount)
				}
			}
		}
	}

	git_stats.Stats(email)
}
