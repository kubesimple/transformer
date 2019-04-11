package transform

// Service is a kubesimple service configured under the `services` section of kubesimple.yml
type Service interface {
	Kind() string
	Validate() error
	Apply() error
}
