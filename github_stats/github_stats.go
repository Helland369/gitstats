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
	for i := daysInLastSixMonths; i > 0; i-- {
		commits[i] = 0
	}

	cal := res.Data.User.ContributionsCollection.ContributionCalendar
	offset := git_stats.Calc_offset()

	for _, w := range cal.Weeks {
		for _, d := range w.ContributionDays {
			t, err := time.Parse("2006-01-02", d.Date)
			if err != nil {
				continue
			}
			daysAgo := git_stats.Count_days_since_date(t) + offset
			if daysAgo == 9999 {
				continue
			}
			if daysAgo > 0 && daysAgo <= daysInLastSixMonths {
				commits[daysAgo] += d.ContributionCount
			}
		}
	}

	return commits
}

