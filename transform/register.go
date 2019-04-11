package transform

import (
	"github.com/kubesimple/transformer/context"
)

const (
	h context.Environment = context.Hosted
	p context.Environment = context.Prod
	d context.Environment = context.Dev
)

func RegisterAllPlugins() HandlerTree {
	r := NewHandlerTree()

	r.Register("app", "1")
	r.Register("postgres", "1", Env(p, d))
	r.Register("disk", "1")
	// Add new plugins here

	return r
}
