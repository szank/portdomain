package persistence

import "github.com/szank/portdomain/pkg/service"

type InMemoryDatabase struct {
	data map[string]service.Port
}

func (d *InMemoryDatabase) Upsert(port service.Port) (bool, error) {
	return false, nil
}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		data: make(map[string]service.Port),
	}
}
