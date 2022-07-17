package service

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
type PortDomainService struct {
	db Database
}

func (s *PortDomainService) Start() error {
	return nil
}

func New(database Database) *PortDomainService {
	return &PortDomainService{
		db: database,
	}
}
