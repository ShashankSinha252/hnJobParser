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
}

func GetComment(id int) (*Comment, error) {

	api := fmt.Sprintf("v0/item/%d.json", id)
	queryURL := baseURL + api

	resp, err := http.Get(queryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed get call: status code: %d", resp.StatusCode)
	}

	comment := &Comment{}
	err = json.NewDecoder(resp.Body).Decode(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (c *Comment) Save() error {
	fileID := "comment/" + strconv.Itoa(c.ID)

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
