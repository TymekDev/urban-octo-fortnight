package main

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
