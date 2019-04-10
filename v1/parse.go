package v1

import (
	"github.com/spf13/viper"
)

// KsDef represents a single kubesimple.yaml app definition
type KsDef struct {
	Version  string
	Appname  string
	Services []Service
}

type Session struct {
	OrgID       string
	UserID      string
	OrgName     string
	UserName    string
	Namespace   string
	Environment string
}

func Parse(v *viper.Viper) error {
	return nil
}
