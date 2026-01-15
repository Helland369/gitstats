package github_stats

import "encoding/json"

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
