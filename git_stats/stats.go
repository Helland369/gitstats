package git_stats

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
	repo, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}
	ref, err := repo.Head()
	if err != nil {
		panic(err)
	}
	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}
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

	for _, k := range keys {
		week := k / 7
		dow := k % 7

		col, ok := cols[week]
		if !ok {
			col = make(column, 7)
		}
		col[dow] = commits[k]
		cols[week] = col
	}
	return cols
}

func print_months() {
	week := get_begining_of_day(time.Now()).Add(-(daysInLastSixMonths * time.Hour * 24))
	month := week.Month()
	fmt.Printf("    ")
	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("    ")
		}
		week = week.Add(7 * time.Hour * 24)
		if week.After(time.Now()) {
			break
		}
 	}
	fmt.Printf("\n")
}

func print_days_col(dow int) {
	switch dow {
	case 0:
	  fmt.Print(" Sun ")
	case 1:
		fmt.Print(" Mon ")
	case 2:
 		fmt.Print(" Tue ")
	case 3:
		fmt.Print(" Wed ")
	case 4:
		fmt.Print(" Thu ")
	case 5:
		fmt.Print(" Fri ")
	case 6:
		fmt.Print(" Sat ")
	default:
		fmt.Print("     ")
	}
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
		fmt.Print(escape + "  - " + "\033[0m")
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

	weeks := weeksInLastSixMonths + 1

	for dow := range 7 {
		print_days_col(dow)

		for w := range weeks {
			if col, ok := cols[w]; ok && len(col) > dow {
				print_cell(col[dow], false)
			} else {
				print_cell(0, false)
			}
		}
		fmt.Print("\n")
	}
}

func print_commit_stats(commits map[int]int) {
	keys := sort_map_into_slices(commits)
	cols := build_cols(keys, commits)
	print_cells(cols)
}

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

func Stats(email string) {
	commits := process_repositories(email)
	print_commit_stats(commits)
}

// export functions
func Calc_offset() int { return calc_offset() }
func Count_days_since_date(t time.Time) int { return count_days_since_date(t) }
func Print_commit_stats(commits map[int]int) { print_commit_stats(commits) }
