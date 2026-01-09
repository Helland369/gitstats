package main

import (
	"fmt"
	"sort"
	"time"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)


const daysInLastSixMonths = 183
const outOfRange = 9999
const weeksInLastSixMonths = 26

type column []int

func get_begining_of_day(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return startOfDay
}

func count_days_since_date(date time.Time) int {
	days := 0
	now := get_begining_of_day(time.Now())
	for date.Before(now) {
		date = date.Add(time.Hour * 24)
		days++
		if days > daysInLastSixMonths {
			return outOfRange
		}
	}
	return days
}

func calc_offset() int {
	var offset int
	weekday := time.Now().Weekday()
	switch weekday {
	case time.Sunday:
		offset = 7
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	case time.Saturday:
		offset = 1
	}
	return offset
}

func fill_commits(email string, path string, commits map[int]int) map[int]int {
	// instance of git repo object from path
	repo, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}
	// get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		panic(err)
	}
	// get the commit history starting from HEAD
	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}
	// iterate the commits
	offset := calc_offset()
	err = iterator.ForEach(func(c *object.Commit) error {
		daysAgo := count_days_since_date(c.Author.When) + offset
		if c.Author.Email != email {
			return nil
		}
		if daysAgo != outOfRange {
			commits[daysAgo]++
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return commits
}

func sort_map_into_slices(m map[int]int) []int {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func build_cols(keys []int, commits map[int]int) map[int]column {
	cols := make(map[int]column)
	col := column{}
	for _, k := range keys {
		week := int(k / 7)
		dayinweek := k % 7
		if dayinweek == 0 {
			col = column{}
		}
		col = append(col, commits[k])
		if dayinweek == 6 {
			cols[week] = col
		}
	}
	return cols
}

func print_months() {
	week := get_begining_of_day(time.Now()).Add(-(daysInLastSixMonths * time.Hour * 24))
	month := week.Month()
	fmt.Printf("      ")
	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("      ")
		}
		week = week.Add(7 * time.Hour * 24)
		if week.After(time.Now()) {
			break
		}
 	}
	fmt.Printf("\n")
}

func print_days_col(day int) {
	out := "   "
	switch day {
	case 1:
		out = " Mon "
	case 2:
		out = " Wed "
	case 3:
		out = " Fri "
	}
	fmt.Printf(out)
}

func print_cell(val int, today bool) {
	escape := "\033[0;37;30m"
	switch {
	case val > 0 && val < 5:
		escape = "\033[1;30;47m"
	case val >= 5 && val < 10:
		escape = "\033[1;30;43m"
	case val >= 10:
		escape = "\033[1;30;42m"
	}
	if today {
		escape = "\033[1;37;45m"
	}
	if val == 0 {
		fmt.Printf(escape + "  - " + "\033[0m")
		return
	}
	str := "  %d "
	switch {
	case val >= 10:
		str = " %d "
	case val >= 100:
		str = "%d "
	}
	fmt.Printf(escape + str + "\033[0m", val)
}

func print_cells(cols map[int]column) {
	print_months()
	for j := 6; j >= 0; j-- {
		for i := weeksInLastSixMonths + 1; i >= 0; i-- {
			if i == weeksInLastSixMonths + 1 {
				print_days_col(j)
			}
			if col, ok := cols[i]; ok {
				if i == 0 && j == calc_offset() -1 {
					print_cell(col[j], true)
					continue
				} else {
					if len(col) > j {
						print_cell(col[j], false)
						continue
					}
				}
			}
			print_cell(0, false)
		}
		fmt.Printf("\n")
	}
}

func print_commit_stats(commits map[int]int) {
	keys := sort_map_into_slices(commits)
	cols := build_cols(keys, commits)
	print_cells(cols)
}

// process_repositories given a user email, returns the commits made in the last 6 months
func process_repositories(email string) map[int]int {
	filePath := get_dot_file_path()
	repos := parse_file_lines_to_slice(filePath)
	daysInMap := daysInLastSixMonths
	commits := make(map[int]int, daysInMap)
	for i := daysInMap; i > 0; i-- {
		commits[i] = 0;
	}
	for _, path := range repos {
		commits = fill_commits(email, path, commits)
	}
	return commits
}

func stats(email string) {
	commits := process_repositories(email)
	print_commit_stats(commits)
}
