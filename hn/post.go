package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Post struct {
	Title          string `json:"title"`
	Content        string `json:"text"`
	BaseCommentIDs []int  `json:"kids"`
	Poster         string `json:"by"`
}

func GetPost(postID int) (*Post, error) {
	api := fmt.Sprintf("v0/item/%d.json", postID)
	queryURL := baseURL + api

	resp, err := http.Get(queryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed get call: status code: %d", resp.StatusCode)
	}

	post := &Post{}
	err = json.NewDecoder(resp.Body).Decode(post)
	if err != nil {
		return nil, err
	}

	return post, nil

}
