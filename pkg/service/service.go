package service

import (
	"errors"
	"fmt"
	"io"
)

// Port struct holds the information about ports (as in physical places for loading unloading ships,
// not unsigned integers, it gets me every time ;))
type Port struct {
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"` // TODO: validate that len == 2 if exists
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"` // TODO: validate if it's a correct string
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code"` // Should be unique, no?
}

// Database is an interface used to communicate with the persistent storage layer.
// Takes in the Port struct, returns boolean specifying if the operation was an upsert (true),
// or insert (false). Errors can be transient or not, database layer should provide an appropriate
// implementation for both cases.
type Database interface {
	Upsert(Port) (bool, error)
}

type RecordProvider interface {
	Next() (Port, error)
}

// PortDomainService contain the main business logic of persisting the incoming port data into a
// persistent storage.
type PortDomainService struct {
	db       Database
	provider RecordProvider
}

// Starts starts the service. Returns an error if the processing pipeline fails. Returns nil if the
// whole pipeline is completed.
func (s *PortDomainService) Start() error {
	for {
		record, err := s.provider.Next()
		switch {
		case err == nil:

		case errors.Is(err, io.EOF):
			// we are done
			return nil
		default:
			return fmt.Errorf("error while retrieving records from the provider: %w", err)
		}

		// we currently don't care if it's an upsert
		_, err = s.db.Upsert(record)
		if err != nil {
			return fmt.Errorf("error inserting the data into the database: %w", err)
		}
	}
}

// New returns a new instance of the PortDomainService
func New(database Database, provider RecordProvider) *PortDomainService {
	return &PortDomainService{
		db:       database,
		provider: provider,
	}
}
