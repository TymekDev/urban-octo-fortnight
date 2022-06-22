package main

type Storage interface {
	NewUser(username string) error
}
