package v1

import (
	"github.com/kubesimple/transformer/context"
	"github.com/spf13/viper"
)

// KsDef represents a single kubesimple.yaml app definition and associated session state
type KsDef struct {
	Version  string
	Appname  string
	Services []Service
}

// Parse reads the kubesimple config file, validates the configuration, and generates a DAG of transformation
// operations to create the resource manifests for the given session environment.
func Parse(v *viper.Viper, s context.Session) error {
	return nil
}
