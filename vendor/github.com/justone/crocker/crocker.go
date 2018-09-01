package crocker

import (
	"github.com/docker/docker-credential-helpers/client"
	"github.com/docker/docker-credential-helpers/credentials"
)

// Crocker is a helper to help with the docker credential helpers.
type Crocker struct {
	Prog client.ProgramFunc
}

// New creates a Crocker instance with the default strategy (StockStrategy).
func New() (*Crocker, error) {
	return NewWithStrategy(StockStrategy{})
}

// New creates a Crocker instance with the supplied strategy.
func NewWithStrategy(strat Strategy) (*Crocker, error) {
	helper, err := strat.Helper()
	if err != nil {
		return nil, err
	}

	return &Crocker{
		client.NewShellProgramFunc(helper),
	}, nil
}

// Erase erases credentials using the detected helper.
func (c *Crocker) Erase(serverURL string) error {
	return client.Erase(c.Prog, serverURL)
}

// Get gets credentials using the detected helper.
func (c *Crocker) Get(serverURL string) (*credentials.Credentials, error) {
	return client.Get(c.Prog, serverURL)
}

// List lists credentials using the detected helper.
func (c *Crocker) List() (map[string]string, error) {
	return client.List(c.Prog)
}

// Store stores credentials using the detected helper.
func (c *Crocker) Store(creds *credentials.Credentials) error {
	return client.Store(c.Prog, creds)
}
