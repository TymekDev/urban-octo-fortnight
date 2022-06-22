package main

type Model map[User]UserData

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
