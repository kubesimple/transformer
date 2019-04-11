package transform

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/kubesimple/transformer/context"
	"github.com/pkg/errors"
)

const (
	_ = iota
	_
	Mi float64 = 1 << (10 * iota)
	Gi
)

type Validation func(s string) error

func Validate(s string, funcs ...Validation) error {
	for _, f := range funcs {
		if err := f(s); err != nil {
			return err
		}
	}
	return nil
}

func Oneof(allowable ...string) Validation {
	return func(s string) error {
		for _, a := range allowable {
			if a == s {
				return nil
			}
		}
		return fmt.Errorf("value %s failed OneOf validation, not contained in %+v", s, allowable)
	}
}

func Required() Validation {
	return Requiredf("")
}

func Requiredf(errfmt string, args ...interface{}) Validation {
	return func(s string) error {
		switch {
		case len(s) == 0:
			return fmt.Errorf(errfmt, args...)
		default:
			return nil
		}
	}
}

func AllowedEnvironment() Validation {
	return func(s string) error {
		env := context.Environment(s)
		switch env {
		case context.Prod, context.Dev, context.Hosted:
			return nil
		default:
			return fmt.Errorf("unknown environment: %s", s)
		}
	}
}

func ValidMemory() Validation {
	return func(s string) error {
		switch {
		case strings.HasSuffix(s, "Mi"), strings.HasSuffix(s, "Gi"):
			return nil
		default:
			return fmt.Errorf("invalid memory limit: %s", s)
		}
	}
}

func MemoryNotExceed(limit string) Validation {
	return func(s string) error {
		lim64, err := memToFloat64(limit)
		if err != nil {
			return errors.Wrapf(err, "could not validate memory limit")
		}
		s64, err := memToFloat64(s)
		if err != nil {
			return errors.Wrapf(err, "could not validate memory limit")
		}
		switch {
		case s64 > lim64:
			return fmt.Errorf("allocated memory of %s exceeds limit of %s", s, limit)
		default:
			return nil
		}
	}
}

// TODO: make public?
func memToFloat64(s string) (float64, error) {
	switch {
	case strings.HasSuffix(s, "Mi"):
		mem, err := strconv.ParseFloat(strings.TrimSuffix(s, "Mi"), 64)
		if err != nil {
			return 0, fmt.Errorf("failed to convert memory string %s to float", s)
		}
		return mem * Mi, nil
	case strings.HasSuffix(s, "Gi"):
		mem, err := strconv.ParseFloat(strings.TrimSuffix(s, "Gi"), 64)
		if err != nil {
			return 0, fmt.Errorf("failed to convert memory string %s to float", s)
		}
		return mem * Gi, nil
	default:
		return 0, fmt.Errorf("failed to convert memory string %s to float, unknown units", s)
	}
}

func float64ToMem(m float64) string {
	return fmt.Sprintf("%dMi", int(math.Ceil(m/Mi)))
}
