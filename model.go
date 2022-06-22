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

var _ Storage = (*Model)(nil)

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
	m.data[user].Run()
	return nil
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
	Factories factories
	Resources *struct {
		Iron   int
		Copper int
		Gold   int
	}
}

func (ud *userData) Run() {
	copper := make(chan int)
	iron := make(chan int)
	gold := make(chan int)
	go ud.Factories.Copper.Run(copper)
	go ud.Factories.Iron.Run(iron)
	go ud.Factories.Gold.Run(gold)
	go func() {
		for {
			ud.Resources.Copper += <-copper
		}
	}()
	go func() {
		for {
			ud.Resources.Iron += <-iron
		}
	}()
	go func() {
		for {
			ud.Resources.Gold += <-gold
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
		Factories: factories{
			Iron:   newFactory(iron),
			Copper: newFactory(copper),
			Gold:   newFactory(gold),
		},
	}
}
