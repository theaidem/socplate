package main

import (
	"fmt"
	"strings"
	"time"
)

type user struct {
	ID      int
	Name    string
	Avatar  string
	Created int64
	Updated int64
}

type userResponse struct {
	ID     int     `msgpack:"id"`
	Name   string  `msgpack:"name"`
	Avatar string  `msgpack:"avatar"`
	Room   *string `msgpack:"room"`
}

func newUser(id int) (*user, error) {

	u, err := generateUser()
	if err != nil {
		return nil, err
	}

	return &user{
		ID:      id,
		Name:    fmt.Sprintf("%s %s", strings.Title(u.Results[0].Name.First), strings.Title(u.Results[0].Name.Last)),
		Avatar:  u.Results[0].Picture.Medium,
		Created: time.Now().Unix(),
		Updated: time.Now().Unix(),
	}, nil
}

func (u *user) userResponse() *userResponse {
	return &userResponse{
		ID:     u.ID,
		Name:   u.Name,
		Avatar: u.Avatar,
	}
}

func (u *userResponse) addRoom(roomID string) *userResponse {
	u.Room = &roomID
	return u
}
