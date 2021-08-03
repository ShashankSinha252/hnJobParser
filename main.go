package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/deadmangareader/hnJobParser/hn"
)

const (
	userID = "whoishiring"
)

type Config struct {
	baseDir string
}

func getConfig() (*Config, error) {
	config := &Config{}

	flag.StringVar(&config.baseDir, "commentDir", "comment", "Directory to store comments")
	flag.Parse()

	// Sanity check on parameters
	stat, err := os.Stat(config.baseDir)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
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

	//Latest post by user
	postID := user.PostIDs[0]
	post, err := hn.GetPost(postID)
	if err != nil {
		fmt.Printf("unable to get post[%d]: %v\n", postID, err)
		return
	}

	if !titleMatch(post.Title) {
		fmt.Println("Post of interest not found")
		return
	}

	fmt.Printf("Post[%s] has %d comments\n", post.Title, len(post.BaseCommentIDs))
	getAndSaveComments(config, post.BaseCommentIDs)
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

func getAndSaveComments(config *Config, commentIDs []int) {
	var wg sync.WaitGroup
	wg.Add(len(commentIDs))

	for _, commentID := range commentIDs {
		go func(id int, w *sync.WaitGroup) {
			defer w.Done()

			comment, err := hn.GetComment(id)
			if err != nil {
				fmt.Printf("Invalid comment[%d]: %v\n", id, err)
				return
			}

			err = comment.Save(config.baseDir)
			if err != nil {
				fmt.Printf("Error in saving comment[%d]: %v\n", id, err)
			}
		}(commentID, &wg)
	}

	wg.Wait()

}
