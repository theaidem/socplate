package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const randUserEndpoint = "https://randomuser.me/api/"

type randomUser struct {
	Results []struct {
		Gender string `json:"gender"`
		Name   struct {
			Title string `json:"title"`
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Email string `json:"email"`
		Login struct {
			UUID     string `json:"uuid"`
			Username string `json:"username"`
		} `json:"login"`
		Phone   string `json:"phone"`
		Picture struct {
			Large     string `json:"large"`
			Medium    string `json:"medium"`
			Thumbnail string `json:"thumbnail"`
		} `json:"picture"`
	} `json:"results"`
}

func generateUser() (*randomUser, error) {
	resp, err := http.Get(randUserEndpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := randomUser{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil

}
