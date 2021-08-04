package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/deadmangareader/hnJobParser/hn"
)

const (
	userID = "whoishiring"
)

const (
	mapFile = "map.txt"
	newLine = "\n"
)

type Config struct {
	baseDir string
}

func getConfig() (*Config, error) {
	config := &Config{}

	flag.StringVar(&config.baseDir, "commentDir", "comment", "Directory to store comments")
	flag.Parse()

	// Sanity check on parameters
	config.baseDir = sanitizeDirPath(config.baseDir)
	if !isDirectory(config.baseDir) {
		return nil, fmt.Errorf("not a directory: %s", config.baseDir)
	}

	return config, nil
}

func main() {

	config, err := getConfig()
	if err != nil {
		fmt.Printf("unable to get a valid config: %v\n", err)
		return
	}

	user, err := hn.GetUser(userID)
	if err != nil {
		fmt.Printf("unable to get user[%s]: %v\n", userID, err)
		return
	}

	post, err := findPost(user)
	if err != nil {
		fmt.Printf("unable to get post: %v\n", err)
		return
	}

	processPost(config.baseDir, post)
}

// Get list of new comments and save their content
func processPost(dir string, post *hn.Post) {
	fmt.Printf("Post[%s] has %d total comments\n",
		post.Title, len(post.BaseCommentIDs))

	oldMap := oldComments(dir)
	idList := newComments(post.BaseCommentIDs, oldMap)
	fmt.Printf("Post[%s] has %d new comments\n",
		post.Title, len(idList))

	ch := saveComment(dir, idList)
	saveMap(dir, ch, oldMap)
}

func findPost(user *hn.User) (*hn.Post, error) {
	postID := user.PostIDs[0]
	post, err := hn.GetPost(postID)
	if err != nil {
		return nil, err
	}

	if !titleMatch(post.Title) {
		return nil, fmt.Errorf("post of interest not found")
	}

	return post, nil
}

// titleMatch checks if we have the
// post of interest by inspecting title
func titleMatch(title string) bool {
	now := time.Now()

	month := now.Month().String()
	if !strings.Contains(title, month) {
		return false
	}

	year := strconv.Itoa(now.Year())
	if !strings.Contains(title, year) {
		return false
	}

	subsTitle := "hiring"
	return strings.Contains(title, subsTitle)
}

func extractComment(id int, w *sync.WaitGroup, ch chan<- int, dir string) {
	// TODO: Move to a worker pool setup for rate limiting
	defer w.Done()

	c, err := hn.GetComment(id)
	if err != nil {
		fmt.Printf("Invalid comment[%d]: %v\n", id, err)
		return
	}

	err = c.Save(dir)
	if err != nil {
		fmt.Printf("Error in saving comment[%d]: %v\n", id, err)
	}

	ch <- id
}

func saveComment(dir string, ids []int) <-chan int {
	size := len(ids)

	var wg sync.WaitGroup
	wg.Add(size)
	ch := make(chan int, size)

	for _, id := range ids {
		go extractComment(id, &wg, ch, dir)
	}

	wg.Wait()
	close(ch)
	return ch
}

// Saves map of comment IDs to file
func saveMap(dir string, ch <-chan int, oldMap map[int]bool) error {
	name := path.Join(dir, mapFile)
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	for id := range ch {
		writeIntToFile(id, file)
	}

	for id := range oldMap {
		writeIntToFile(id, file)
	}

	return nil
}

// Writes an ID to given file handle
func writeIntToFile(val int, file *os.File) {
	content := strconv.Itoa(val) + newLine
	file.Write([]byte(content))

}

// Gives a map containing old comments, dumped
// from a file
func oldComments(dir string) map[int]bool {
	m := make(map[int]bool)

	file := path.Join(dir, mapFile)
	data, err := os.ReadFile(file)
	if err != nil {
		return m
	}

	getIntsFromString(string(data), m)
	return m
}

// Splits str into lines, then converts each line to
// int and adds it to map
func getIntsFromString(str string, m map[int]bool) {
	lines := strings.Split(str, newLine)
	for _, line := range lines {
		val, err := strconv.Atoi(line)
		if err != nil {
			continue
		}
		m[val] = true
	}
}

// Filter outs old comments from given list of comments
func newComments(all []int, old map[int]bool) []int {
	want := []int{}

	for _, val := range all {
		if old[val] {
			continue
		}
		want = append(want, val)
	}

	return want
}

// Checks if dir is a directory
func isDirectory(dir string) bool {
	stat, err := os.Stat(dir)
	if err != nil {
		return false
	}
	if !stat.IsDir() {
		return false
	}
	return true
}

// sanitizeDirPath removes trailing '/' from dir
// for usage with path.Join()
func sanitizeDirPath(dir string) string {
	lastIndex := len(dir) - 1
	if dir[lastIndex] == '/' {
		return dir[:len(dir)-1]
	}
	return dir
}
