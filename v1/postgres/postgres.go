package postgres

import (
	"github.com/kubesimple/transformer/build"
)

var _ build.Service = PgService{}

// PgService is the definition of a postgres DB service in the kubesimple.yaml file
type PgService struct {
	Image    string
	Memory   string
	Migrate  string
	Init     string
	Replicas int
}

// Kind fulfills the service interface
func (s PgService) Kind() string {
	return Postgres
}
