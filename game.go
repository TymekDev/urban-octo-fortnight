package main

import "time"

type Game interface {
	NewUser(username string) error
	GetUserData(username string) (UserData, error)
	UpgradeFactory(username string, factoryType FactoryType) error
}

type UserData interface {
	Resources() Resources
	Factories() Factories
}

type Resources struct {
	Iron   int
	Copper int
	Gold   int
}

type Factories struct {
	Iron   Factory
	Copper Factory
	Gold   Factory
}

// This is public facing Factory. Ideally, It should be an interface that factory (the private one) fulfills.
type Factory struct {
	Level             int
	Yield             int
	UpgradeInProgress bool
	UpgradeTimeLeft   time.Duration
	UpgradeCost       Resources
}

type FactoryType int

const (
	Iron FactoryType = iota + 1
	Copper
	Gold
)
