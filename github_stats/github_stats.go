package github_stats

import (
	"encoding/json"
	"time"

	"github.com/Helland369/gitstats/git_stats"
)

type ContributionDay struct {
	Date              string `json:"date"`
	Weekday           int    `json:"weekday"`
	ContributionCount int    `json:"contributionCount"`
	Color             string `json:"color"`
}

type Weeks struct {
	FirstDay         string            `json:"firstDay"`
	ContributionDays []ContributionDay `json:"contributionDays"`
}

type ContributionCalendar struct {
	TotalContributions int     `json:"totalContributions"`
	Weeks              []Weeks `json:"weeks"`
}

type ContributionsCollection struct {
	ContributionCalendar ContributionCalendar `json:"contributionCalendar"`
}

type User struct {
	ContributionsCollection ContributionsCollection `json:"contributionsCollection"`
}

type Data struct {
	User User `json:"user"`
}

type Response struct {
	Data Data `json:"data"`
}

func Github_stats(userName string) (Response, error){
	res, err := Get_github_contrib(userName)
	if err != nil {
		return Response{}, err
	}
	
	var response Response
  if err := json.Unmarshal(res, &response); err != nil {
		return Response{}, err
	}

	return response, nil
}

func To_commit_map(res Response) map[int]int {
	const daysInLastSixMonths = 183
	
	commits := make(map[int]int, daysInLastSixMonths)
	for i := 0; i > daysInLastSixMonths; i++ {
		commits[i] = 0
	}

	cal := res.Data.User.ContributionsCollection.ContributionCalendar

	for _, w := range cal.Weeks {
		for _, d := range w.ContributionDays {
			t, err := time.ParseInLocation("2006-01-02", d.Date, time.Local)
			if err != nil {
				continue
			}
			if idx, ok := git_stats.Day_index(t); ok {
				commits[idx] += d.ContributionCount
			}
		}
	}
	
	return commits
}

