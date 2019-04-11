package transform

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/kubesimple/transformer/context"
)

// PluginType is an enum for the types of plugins
type PluginType int

const (
	// Restricted plugins receive only the definition for their particular resource
	Restricted PluginType = iota
	// Unrestricted plugins receive the definition for the entire app
	Unrestricted
)

// HandlerTree holds references to all registered plugins set via the Register method
// and finds an appropriate handler based on kind, version, and environment
type HandlerTree interface {
	Register(kind string, version string, opts ...RegistrationOpt)
	FindKindHandler(kind string, version string, env context.Environment) (Handler, error)
	Validate(e context.Environment) error
}

type handlerTree struct {
	handlers map[string]Handler
}

// Handler represents an executable handler for a kind in one or more environments
type Handler struct {
	Kind       string
	Exe        string
	PluginType PluginType
	Version    string
	Env        []context.Environment
}

// NewHandlerTree returns a HandlerTree interface ready for registering plugins
func NewHandlerTree() HandlerTree {
	return &handlerTree{
		handlers: make(map[string]Handler),
	}
}

// RegistrationOpt passes an option to configure the handler instead of the defaults
type RegistrationOpt func(h *Handler)

// Exe sets the executable name to something other than the default value `ks-plugin-<kind>-v<version>` (e.g., ks-plugin-app-v1)
func Exe(n string) RegistrationOpt {
	return func(h *Handler) {
		h.Exe = n
	}
}

// Env sets the handler for environments other than the default of all environments
func Env(env ...context.Environment) RegistrationOpt {
	return func(h *Handler) {
		h.Env = env
	}
}

// Unrestrict sets this handler to operate in unrestricted mode (default: restricted)
func Unrestrict() RegistrationOpt {
	return func(h *Handler) {
		h.PluginType = Unrestricted
	}
}

// Register sets the executable handler for the kind of resource based on the version.  Register assumes the executable is
// named as `ks-plugin-<kind>-v<version>` (e.g., ks-plugin-app-v1), restricted operation, and for all environments.  To change
// these defaults, use the appropriate registration option (Exe, Env, Unrestrict)
func (h *handlerTree) Register(kind string, version string, opts ...RegistrationOpt) {
	// TODO: panic if the handler does not exist in the path or no environments specified, which represents a fatal unrecoverable
	// error in transformer initialization
	handler := Handler{
		Kind:       kind,
		Version:    version,
		Exe:        fmt.Sprintf("ks-plugin-%s-v%s", kind, version),
		PluginType: Restricted,
		Env:        []context.Environment{context.Prod, context.Hosted, context.Dev},
	}
	for _, opt := range opts {
		opt(&handler)
	}
	for _, e := range handler.Env {
		h.handlers[strings.Join([]string{kind, version, string(e)}, ":")] = handler
	}
}

// FindKindHandler returns the plugin handler based on the kind of resource, plugin version, and environment
func (h *handlerTree) FindKindHandler(kind string, version string, env context.Environment) (Handler, error) {
	handler, ok := h.handlers[strings.Join([]string{kind, version, string(env)}, ":")]
	if !ok {
		return Handler{}, fmt.Errorf("no handler found for kind=%s version=%s env=%s", kind, version, env)
	}
	return handler, nil
}

// Validate ensures that every plugin exists on the path and for the environment specified
func (h *handlerTree) Validate(e context.Environment) error {

	// check if handler is defined for environment
	contains := func(envs []context.Environment) bool {
		for _, env := range envs {
			if env == e {
				return true
			}
		}
		return false
	}

	// return error if not in path and handler defined for this environment
	for _, handler := range h.handlers {
		if _, err := exec.LookPath(handler.Exe); err != nil && contains(handler.Env) {
			return fmt.Errorf("plugin for kind %s not found on path: %s", handler.Kind, handler.Exe)
		}
	}
	return nil
}
