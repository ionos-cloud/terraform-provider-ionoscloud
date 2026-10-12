package failover

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"slices"
	"strings"

	boff "github.com/cenkalti/backoff/v5"
)

// httpMethods lists the HTTP methods accepted in Options.RetryableMethods.
var httpMethods = []string{
	http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut,
	http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace,
}

func isHTTPMethod(method string) bool {
	return slices.Contains(httpMethods, strings.ToUpper(strings.TrimSpace(method)))
}

// Validate checks the failover options. It returns nil when the receiver is nil or
// when every value is acceptable, otherwise an aggregated error describing each
// problem found.
func (o *Options) Validate() error {
	if o == nil {
		return nil
	}

	var problems []error

	switch NormalizeStrategy(o.Strategy) {
	case NormalizeStrategy(None),
		NormalizeStrategy(RoundRobin),
		"":
	default:
		problems = append(problems, fmt.Errorf(
			"invalid failover strategy %q, supported values are %q, %q or an empty value",
			o.Strategy, None, RoundRobin,
		))
	}

	// Every retryableMethods entry must be a valid HTTP method.
	for _, method := range o.RetryableMethods {
		if !isHTTPMethod(method) {
			problems = append(problems, fmt.Errorf(
				"invalid failover retryableMethods entry %q, supported values are %s",
				method, strings.Join(httpMethods, ", "),
			))
		}
	}

	if o.MaxRetries < 0 {
		problems = append(problems, fmt.Errorf("failover maxRetries must be >= 0, got %d", o.MaxRetries))
	}

	if b := o.ExponentialBackoff; b != nil {
		if b.InitialInterval < 0 {
			problems = append(problems, fmt.Errorf("failover exponentialBackoff.initialInterval must be >= 0, got %s", b.InitialInterval))
		}
		if b.MaxInterval < 0 {
			problems = append(problems, fmt.Errorf("failover exponentialBackoff.maxInterval must be >= 0, got %s", b.MaxInterval))
		}
		if b.Multiplier != nil && (math.IsNaN(*b.Multiplier) || math.IsInf(*b.Multiplier, 0) || *b.Multiplier <= 0) {
			problems = append(problems, fmt.Errorf("failover exponentialBackoff.multiplier must be > 0, got %v", *b.Multiplier))
		}
		if b.RandomizationFactor != nil && (math.IsNaN(*b.RandomizationFactor) || *b.RandomizationFactor < 0 || *b.RandomizationFactor > 1) {
			problems = append(problems, fmt.Errorf("failover exponentialBackoff.randomizationFactor must be between 0 and 1, got %v", *b.RandomizationFactor))
		}
		// Resolve zero values to the library defaults and make sure the initial interval does not exceed the maximum.
		if b.InitialInterval >= 0 && b.MaxInterval >= 0 {
			initialInterval := b.InitialInterval
			if initialInterval == 0 {
				initialInterval = boff.DefaultInitialInterval
			}
			maxInterval := b.MaxInterval
			if maxInterval == 0 {
				maxInterval = boff.DefaultMaxInterval
			}
			if initialInterval > maxInterval {
				problems = append(problems, fmt.Errorf("failover exponentialBackoff.initialInterval (%s) must not be greater than maxInterval (%s)", initialInterval, maxInterval))
			}
		}
	}

	return errors.Join(problems...)
}
