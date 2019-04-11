package disk

import (
	"github.com/kubesimple/transformer/build"
)

var _ build.Service = DiskService{}

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
