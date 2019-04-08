package build

import (
	"bytes"
	"fmt"
	"os"
	"path"

	v1 "github.com/kubesimple/transformer/build/v1"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// AppDef represents a single kubesimple.yaml app definition
type AppDef struct {
	Version  string
	Appname  string
	Services []Service

	orgID    string
	userID   string
	orgName  string
	userName string
}

func Read() (*AppDef, error) {
	return read(nil)
}

func read(b []byte) (*AppDef, error) {
	v := viper.New()
	setDefaults(v)
	switch {
	case b == nil:
		if err := v.ReadInConfig(); err != nil {
			return nil, errors.Wrap(err, "failed to read kubesimple configuration file")
		}
	default:
		if err := v.ReadConfig(bytes.NewBuffer(b)); err != nil {
			return nil, errors.Wrap(err, "failed to read kubesimple configuration file")
		}
	}
	def := new(AppDef)
	def.Version = v.GetString("version")

	switch def.Version {
	case "1":
		return v1.Parse(def, v)
	default:
		return nil, errors.New(fmt.Sprintf("config: unknown version %s", def.Version))
	}

}

func setDefaults(v *viper.Viper) {
	v.AddConfigPath(".")
	v.SetConfigFile("kubesimple")
	defaults := map[string]string{
		"version": "1",
	}
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	return
}

func ReadDefault() (*AppDef, error) {
	wd, err := os.Getwd()
	if err != nil {
		wd = path.Clean("./")
	}
	return Read(wd)
}
