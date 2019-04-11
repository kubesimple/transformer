package transform

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/gobuffalo/packr/v2"
	"github.com/kubesimple/transformer/context"
	"github.com/pkg/errors"
)

const (
	h context.Environment = context.Hosted
	p context.Environment = context.Prod
	d context.Environment = context.Dev
)

// author table from plugins.toml
type author struct {
	github    string
	publicKey string
}

// plugin table from plugins.toml
type plugin struct {
	kind         string
	majorVersion string
	exe          string
	author       string
	environment  []string
	signature    string
}

// plugins.toml structure
type pluginsTable struct {
	authors []author
	plugins []plugin
}

// RegisterAllPlugins registers every available plugin from plugins.toml. After registering,
// the HandlerTree can be verified to ensure that all plugins for that environment are available and signed.
func RegisterAllPlugins() (HandlerTree, error) {
	r := NewHandlerTree()
	box := packr.New("plugins", "./plugins.toml")
	pluginsTOML, err := box.FindString("plugins.toml")
	if err != nil {
		return nil, errors.Wrap(err, "could not open plugins.toml to register plugins")
	}
	table := new(pluginsTable)
	if _, err := toml.Decode(pluginsTOML, table); err != nil {
		return nil, errors.Wrap(err, "error parsing plugins.toml")
	}
	for _, plugin := range table.plugins {
		opts := []RegistrationOpt{}
		if plugin.majorVersion == "" {
			return nil, fmt.Errorf("major version for kind %s not defined", plugin.kind)
		}
		if len(plugin.environment) > 0 {
			envs := []context.Environment{}
			for _, env := range plugin.environment {
				envs = append(envs, context.Environment(env))
			}
			opts = append(opts, Env(envs...))
		}
		if plugin.exe != "" {
			opts = append(opts, Exe(plugin.exe))
		}
		r.Register(plugin.kind, plugin.majorVersion, opts...)
	}

	return r, nil
}
