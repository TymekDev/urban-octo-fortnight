package main

import "time"

type factoryType int

const (
	iron = iota + 1
	copper
	gold
)

type factory struct {
	Level int
	Type  factoryType
}

func newFactory(facType factoryType) *factory {
	return &factory{
		Level: 1,
		Type:  facType,
	}
}

type factoryMeta struct {
	ProductionInterval time.Duration
	Upgradeable        bool
	UpgradeTime        time.Duration
	UpgradeCost        upgradeCost
}

type upgradeCost struct {
	Iron   int
	Copper int
	Gold   int
}

var _factoryMeta = map[factoryType]map[int]factoryMeta{
	iron: {
		1: {10 * time.Second, true, 15 * time.Second, upgradeCost{300, 100, 1}},
		2: {20 * time.Second, true, 30 * time.Second, upgradeCost{800, 250, 2}},
		3: {40 * time.Second, true, 60 * time.Second, upgradeCost{1600, 500, 4}},
		4: {80 * time.Second, true, 90 * time.Second, upgradeCost{3000, 1000, 8}},
		5: {150 * time.Second, false, 120 * time.Second, upgradeCost{}},
	},
	copper: {
		1: {3 * time.Second, true, 15 * time.Second, upgradeCost{200, 70, 0}},
		2: {7 * time.Second, true, 30 * time.Second, upgradeCost{400, 150, 0}},
		3: {14 * time.Second, true, 60 * time.Second, upgradeCost{800, 300, 0}},
		4: {30 * time.Second, true, 90 * time.Second, upgradeCost{1600, 600, 0}},
		5: {60 * time.Second, false, 120 * time.Second, upgradeCost{}},
	},
	gold: {
		1: {2 * time.Minute, true, 15 * time.Second, upgradeCost{0, 100, 2}},
		2: {3 * time.Minute, true, 30 * time.Second, upgradeCost{0, 200, 4}},
		3: {4 * time.Minute, true, 60 * time.Second, upgradeCost{0, 400, 8}},
		4: {6 * time.Minute, true, 90 * time.Second, upgradeCost{0, 800, 16}},
		5: {8 * time.Minute, false, 120 * time.Second, upgradeCost{}},
	},
}
