package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Name    string `json:"id"`
	PostIDs []int  `json:"submitted"`
}

const (
	baseURL = "https://hacker-news.firebaseio.com/"
)

func GetUser(userID string) (*User, error) {

	api := fmt.Sprintf("v0/user/%s.json", userID)
	queryURL := baseURL + api

	resp, err := http.Get(queryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed get call: status code: %d", resp.StatusCode)
	}

	user := &User{}
	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil

}
