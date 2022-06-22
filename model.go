package main

import (
	"encoding/json"
	"io/ioutil"
)

type Model map[User]UserData

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

type User string

func NewUser(username string) User {
	return User(username)
}

type UserData struct {
	Iron   Factory
	Copper Factory
	Gold   Factory
}

func NewUserData() UserData {
	return UserData{
		Iron:   NewFactory(),
		Copper: NewFactory(),
		Gold:   NewFactory(),
	}
}

type Factory struct {
	Level int
}

func NewFactory() Factory {
	return Factory{
		Level: 1,
	}
}
