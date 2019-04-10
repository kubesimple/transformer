package app

import (
	"github.com/kubesimple/transformer/build"
)

var _ build.Service = AppService{}

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
