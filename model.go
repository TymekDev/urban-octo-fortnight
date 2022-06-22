package main

import (
	"encoding/json"
	"io/ioutil"
)

type Model map[user]userData

func NewModel(path string) (*Model, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Model
	if err := json.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

type user string

func newUser(username string) user {
	return user(username)
}

type userData struct {
	Iron   factory
	Copper factory
	Gold   factory
}

func newUserData() userData {
	return userData{
		Iron:   newFactory(),
		Copper: newFactory(),
		Gold:   newFactory(),
	}
}

type factory struct {
	Level int
}

func newFactory() factory {
	return factory{
		Level: 1,
	}
}
