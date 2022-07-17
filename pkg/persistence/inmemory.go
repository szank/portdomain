package persistence

import "github.com/szank/portdomain/pkg/service"

// InMemoryDatabase stores the port information in the RAM. The data is lost
// when the binary exists. This is a dummy implmementation of the persistence layer.
type InMemoryDatabase struct {
	data map[string]service.Port
}

// Upsert inserts or updated data in the database. Returns true if the
// data with the unique key (the port name in this implementation) has been
// inserted previous and is updated, and false if the object was not found in the db.
func (d *InMemoryDatabase) Upsert(port service.Port) (bool, error) {
	upsert := false
	_, upsert = d.data[port.Name]

	d.data[port.Name] = port

	return upsert, nil
}

// NewInMemoryDatabase returns a new instance of the InMemoryDatabase object.
func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		data: make(map[string]service.Port),
	}
}
