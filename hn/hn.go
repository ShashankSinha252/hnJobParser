package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Comment struct {
	ID      int    `json:"id"`
	Content string `json:"text"`
	Poster  string `json:"by"`
	Deleted bool   `json:"deleted"`
}

type Post struct {
	Title          string `json:"title"`
	Content        string `json:"text"`
	BaseCommentIDs []int  `json:"kids"`
	Poster         string `json:"by"`
}

type User struct {
	Name    string `json:"id"`
	PostIDs []int  `json:"submitted"`
}

const (
	baseURL = "https://hacker-news.firebaseio.com/"
)

const (
	keyPost    = "post"
	keyComment = "comment"
	keyUser    = "user"
)

var API = map[string]string{
	keyPost:    "v0/item/%s.json",
	keyComment: "v0/item/%s.json",
	keyUser:    "v0/user/%s.json",
}

func getEndpoint(key, param string) string {
	APIFmtStr := API[key]
	endpoint := fmt.Sprintf(APIFmtStr, param)
	return baseURL + endpoint

}

func getElement(elem interface{}, param string) error {
	var key string

	switch t := elem.(type) {
	case *Post:
		key = keyPost
	case *Comment:
		key = keyComment
	case *User:
		key = keyUser
	default:
		return fmt.Errorf("unsupported data type: %v", t)
	}

	queryURL := getEndpoint(key, param)

	resp, err := http.Get(queryURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed get call: status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(elem)
	if err != nil {
		return err
	}

	return nil
}

func GetComment(commentID int) (*Comment, error) {

	comment := &Comment{}
	err := getElement(comment, strconv.Itoa(commentID))
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (c *Comment) Save(basedir string) error {
	if c.ID == 0 {
		return nil
	}

	if c.Deleted {
		return fmt.Errorf("deleted comment")
	}

	if len(basedir) == 0 {
		return fmt.Errorf("no basedir provided")
	}

	if basedir[len(basedir)-1] != '/' {
		basedir += "/"
	}
	fileID := basedir + "commment-" + strconv.Itoa(c.ID)

	file, err := os.Create(fileID)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(c.Content))
	if err != nil {
		return err
	}

	return nil
}

func GetPost(postID int) (*Post, error) {
	post := &Post{}
	err := getElement(post, strconv.Itoa(postID))
	if err != nil {
		return nil, err
	}

	return post, nil
}

func GetUser(userID string) (*User, error) {
	user := &User{}
	err := getElement(user, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
