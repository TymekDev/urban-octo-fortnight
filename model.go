package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"
)

type Model struct {
	sync.Mutex
	data map[user]*userData
	path string
}

var _ Game = (*Model)(nil)

func NewEmptyModel(path string) *Model {
	return &Model{
		data: map[user]*userData{},
		path: path,
	}
}

func NewModel(path string) (*Model, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	m := NewEmptyModel(path)
	if err := json.Unmarshal(bytes, &m.data); err != nil {
		return nil, err
	}
	for _, userData := range m.data {
		go userData.Run()
	}
	return m, nil
}

func (m *Model) NewUser(username string) error {
	user := newUser(username)
	if _, ok := m.data[user]; ok {
		return errors.New("user already exists")
	}
	m.data[user] = newUserData()
	m.save() // FIXME: ideally this shouldn't be synchronous
	go m.data[user].Run()
	return nil
}

func (m *Model) GetUserData(username string) (UserData, error) {
	user := newUser(username)
	userData, ok := m.data[user]
	if !ok {
		return nil, errors.New("user does not exist")
	}
	return userData, nil
}

func (m *Model) UpgradeFactory(username string, factoryType FactoryType) error {
	// TODO: remove duplication without extending UserData interface with UpgradeFactory method
	user := newUser(username)
	userData, ok := m.data[user]
	if !ok {
		return errors.New("user does not exist")
	}
	return userData.UpgradeFactory(factoryType)
}

func (m *Model) save() error {
	m.Lock()
	defer m.Unlock()
	bytes, err := json.Marshal(m.data)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(m.path, bytes, 0644); err != nil {
		return err
	}
	return nil
}

type user string

func newUser(username string) user {
	return user(username)
}

type userData struct {
	UserFactories factories
	UserResources *Resources
}

var _ UserData = (*userData)(nil)

func (ud *userData) Resources() Resources {
	return *ud.UserResources
}

func (ud *userData) Factories() Factories {
	return Factories{
		Iron:   ud.UserFactories.Iron.ToFactory(),
		Copper: ud.UserFactories.Copper.ToFactory(),
		Gold:   ud.UserFactories.Gold.ToFactory(),
	}
}

func (ud *userData) Run() {
	copper := make(chan int)
	iron := make(chan int)
	gold := make(chan int)
	go ud.UserFactories.Copper.Run(copper)
	go ud.UserFactories.Iron.Run(iron)
	go ud.UserFactories.Gold.Run(gold)
	go func() {
		for {
			ud.UserResources.Copper += <-copper
		}
	}()
	go func() {
		for {
			ud.UserResources.Iron += <-iron
		}
	}()
	go func() {
		for {
			ud.UserResources.Gold += <-gold
		}
	}()
}

func (ud *userData) UpgradeFactory(factoryType FactoryType) error {
	var f *factory
	switch factoryType {
	case Iron:
		f = ud.UserFactories.Iron
	case Copper:
		f = ud.UserFactories.Copper
	case Gold:
		f = ud.UserFactories.Gold
	default:
		return errors.New("unknown factory type")
	}
	cost := f.Meta.UpgradeCost
	if cost.Iron > ud.UserResources.Iron || cost.Copper > ud.UserResources.Copper || cost.Gold > ud.UserResources.Gold {
		return errors.New("insufficient resources")
	}
	go f.Upgrade()
	return nil
}

type factories struct {
	Iron   *factory
	Copper *factory
	Gold   *factory
}

func newUserData() *userData {
	return &userData{
		UserFactories: factories{
			Iron:   newFactory(Iron),
			Copper: newFactory(Copper),
			Gold:   newFactory(Gold),
		},
		UserResources: &Resources{},
	}
}
