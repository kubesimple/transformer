package transformer

import (
	"bytes"
	"fmt"

	v1 "github.com/kubesimple/transformer/v1"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Transform() error {
	return transform(nil)
}

func transform(b []byte) error {
	v := viper.New()
	setDefaults(v)
	switch {
	case b == nil:
		if err := v.ReadInConfig(); err != nil {
			log.Errorf("failed to read kubesimple configuration file: %s", err)
			return errors.Wrap(err, "failed to read kubesimple configuration file")
		}
	default:
		if err := v.ReadConfig(bytes.NewBuffer(b)); err != nil {
			log.Errorf("failed to read kubesimple configuration file: %s", err)
			return errors.Wrap(err, "failed to read kubesimple configuration file")
		}
	}
	version := v.GetString("version")

	switch version {
	case "1":
		return v1.Parse(v)
	default:
		return errors.New(fmt.Sprintf("config: unknown version %s", version))
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
}
