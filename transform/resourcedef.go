package v1

const (
	App      = "app"
	Disk     = "disk"
	Postgres = "postgres"
)

// Service is a kubesimple service configured under the `services` section of kubesimple.yml
type Service interface {
	Kind() string
	Validate() error
	Apply() error
}
