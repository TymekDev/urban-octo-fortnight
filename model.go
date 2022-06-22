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

type factories struct {
	Iron   *factory
	Copper *factory
	Gold   *factory
}

func newUserData() *userData {
	return &userData{
		UserFactories: factories{
			Iron:   newFactory(iron),
			Copper: newFactory(copper),
			Gold:   newFactory(gold),
		},
		UserResources: &Resources{},
	}
}
