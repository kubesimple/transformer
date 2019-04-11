package context

import (
	"fmt"
	"strings"
)

// Environment is an enum to represent the deployment environment
type Environment string

const (
	// Prod is a production deployment on a standard K8S cluster.  Prod deployments may be missing
	// some features of the KS hosted environment that must be installed manually
	Prod Environment = "prod"
	// Hosted is a production deployment specific to the Kubesimple architecture
	Hosted Environment = "hosted"
	// Dev is a dev deployment for a user's own development machine, not suitable for production
	Dev Environment = "dev"
)

// Session represents details of a user's KS account and the deployment environment.  When a user
// is not logged in or does not have a KS account, the required values will be generated in dev mode
// or optionally set via env.
type Session struct {
	OrgID       string
	UserID      string
	OrgName     string
	UserName    string
	Namespace   string
	Environment Environment
}

// GetEnvironment parses the environment string determined from command line arguments and env variables
// and validates that it matches one of the allowed transform/build environments
func GetEnvironment(e string) (Environment, error) {
	switch env := strings.ToLower(e); env {
	case "prod", "production":
		return Prod, nil
	case "hosted":
		return Hosted, nil
	case "dev", "development":
		return Dev, nil
	default:
		return "", fmt.Errorf("unknown build environment: %s", env)
	}
}
