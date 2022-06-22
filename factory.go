package main

import (
	"encoding/json"
	"time"
)

type factory struct {
	Level             int
	Type              FactoryType
	UpgradeInProgress bool
	UpgradeEndTime    time.Time
	Meta              *factoryMeta `json:"-"`
}

func newFactory(facType FactoryType) *factory {
	return &factory{
		Level: 1,
		Type:  facType,
		Meta:  _factoryMeta[facType][1],
	}
}

func (f *factory) Run(c chan int) {
	for range time.Tick(f.Meta.YieldInterval) {
		c <- f.Meta.Yield
	}
}

func (f *factory) ToFactory() Factory {
	result := Factory{
		Level:             f.Level,
		Yield:             f.Meta.Yield,
		UpgradeInProgress: f.UpgradeInProgress,
	}
	if f.UpgradeInProgress {
		result.UpgradeTimeLeft = f.UpgradeEndTime.Sub(time.Now())
	} else {
		result.UpgradeCost = f.Meta.UpgradeCost
	}
	return result
}

func (f *factory) Upgrade() {
	f.UpgradeInProgress = true
	f.UpgradeEndTime = time.Now().Add(f.Meta.UpgradeTime)
	time.Sleep(f.Meta.UpgradeTime) // FIXME: sleep seems to be an anti-pattern
	f.Level++
	f.UpgradeInProgress = false
	f.Meta = _factoryMeta[f.Type][f.Level]
}

func (f *factory) UnmarshalJSON(data []byte) error {
	var f2 struct {
		Level int
		Type  FactoryType
	}
	if err := json.Unmarshal(data, &f2); err != nil {
		return err
	}
	f.Level = f2.Level
	f.Type = f2.Type
	f.Meta = _factoryMeta[f2.Type][f2.Level]
	return nil
}

type factoryMeta struct {
	Yield         int
	YieldInterval time.Duration
	Upgradeable   bool
	UpgradeTime   time.Duration
	UpgradeCost   Resources
}

var _factoryMeta = map[FactoryType]map[int]*factoryMeta{
	Iron: {
		1: {10, time.Second, true, 15 * time.Second, Resources{300, 100, 1}},
		2: {20, time.Second, true, 30 * time.Second, Resources{800, 250, 2}},
		3: {40, time.Second, true, 60 * time.Second, Resources{1600, 500, 4}},
		4: {80, time.Second, true, 90 * time.Second, Resources{3000, 1000, 8}},
		5: {150, time.Second, false, 120 * time.Second, Resources{}},
	},
	Copper: {
		1: {3, time.Second, true, 15 * time.Second, Resources{200, 70, 0}},
		2: {7, time.Second, true, 30 * time.Second, Resources{400, 150, 0}},
		3: {14, time.Second, true, 60 * time.Second, Resources{800, 300, 0}},
		4: {30, time.Second, true, 90 * time.Second, Resources{1600, 600, 0}},
		5: {60, time.Second, false, 120 * time.Second, Resources{}},
	},
	Gold: {
		1: {2, time.Minute, true, 15 * time.Second, Resources{0, 100, 2}},
		2: {3, time.Minute, true, 30 * time.Second, Resources{0, 200, 4}},
		3: {4, time.Minute, true, 60 * time.Second, Resources{0, 400, 8}},
		4: {6, time.Minute, true, 90 * time.Second, Resources{0, 800, 16}},
		5: {8, time.Minute, false, 120 * time.Second, Resources{}},
	},
}
