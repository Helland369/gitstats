package git_stats

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"
)

// scan for .git directories
func scan_git_directories(dirs []string, dir string) []string {
	dir = strings.TrimSuffix(dir, "/")
	f, err := os.Open(dir)
	if err != nil {
		log.Fatal(err)
	}

	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	var path string
	for _, file := range files {
		if file.IsDir() {
			path = dir + "/" + file.Name()
			if file.Name() == ".git"{
				path = strings.TrimSuffix(path, "/.git")
				fmt.Println(path)
				dirs = append(dirs, path)
				continue
			}
			if file.Name() == "vendor" || file.Name() == "node_modules" {
				continue
			}
			dirs = scan_git_directories(dirs, path)
		}
	}

	return dirs
}

func recursive_scan_directories(dir string) []string {
	return scan_git_directories(make([]string, 0), dir)
}

func get_dot_file_path() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dotFile := usr.HomeDir + "/.gitlocalstats"
	return dotFile
}

func open_file(filePath string) *os.File {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			_, err := os.Create(filePath)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	return f
}

func parse_file_lines_to_slice(filePath string) []string {
	f := open_file(filePath)
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			panic(err)
		}
	}
	return lines
}

func slice_contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func join_slices(new []string, existing []string) []string {
	for _, i := range new {
		if !slice_contains(existing, i) {
			existing = append(existing, i)
		}
	}
	return existing
}

func dump_string_slice_to_file(repos []string, filePath string) {
	content := strings.Join(repos, "\n")
	os.WriteFile(filePath, []byte(content), 0755)
}

func add_new_slice_elements_to_file(filePath string, newRepos []string) {
	existingRepos := parse_file_lines_to_slice(filePath)
	repos := join_slices(newRepos, existingRepos)
	dump_string_slice_to_file(repos, filePath)
}

// scan a new directory for a Git repository
func Scan(dir string) {
	fmt.Printf("Found Directories:\n\n")
	repositories := recursive_scan_directories(dir)
	filepath := get_dot_file_path()
	add_new_slice_elements_to_file(filepath, repositories)
	fmt.Printf("\n\nSuccessfully added\n\n")
}

