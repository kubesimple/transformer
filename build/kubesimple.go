package build

const (
	App      = "app"
	Database = "db"
	Disk     = "disk"
)

var _ Service = AppService{}
var _ Service = DbService{}
var _ Service = DiskService{}

// Service is a kubesimple service configured under the `services` section of kubesimple.yml
type Service interface {
	Kind() string
}

// AppService is the definition of a app using a custom image or built using a buildpack in the repo root
type AppService struct {
	Image      string
	Build      string
	Public     bool
	PublicOpts PubOpt
	Memory     string
	Replicas   int
}

// PubOpt are options for publically exposing an app
type PubOpt struct {
	Host string
	Path string
}

// Kind fulfills the service interface
func (s AppService) Kind() string {
	return App
}

// DbService is the definition of a database service in the kubesimple.yaml file
type DbService struct {
	Image    string
	Memory   string
	Migrate  string
	Init     string
	Replicas int
}

// Kind fulfills the service interface
func (s DbService) Kind() string {
	return Database
}

// DiskService is the definition of a abstracted disk service in the kubesimple.yaml file
type DiskService struct {
	Size       string
	Path       string
	Restricted bool
}

// Kind fulfills the service interface
func (s DiskService) Kind() string {
	return Disk
}
