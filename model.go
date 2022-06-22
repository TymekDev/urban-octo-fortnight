package main

type Model map[User]UserData

type User string

type UserData struct {
	Iron   Factory
	Copper Factory
	Gold   Factory
}

type Factory struct {
	Level int
}
