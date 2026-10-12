package fileconfiguration

import (
	"errors"
	"fmt"
	"strings"
)

// knownProducts is the set of product names accepted in a file config. It is built
// from the product-name constants above; new products must be added here too.
var knownProducts = map[string]struct{}{
	Autoscaling:             {},
	APIGateway:              {},
	CDN:                     {},
	Cert:                    {},
	Cloud:                   {},
	ContainerRegistry:       {},
	DNS:                     {},
	Mongo:                   {},
	ObjectStorageManagement: {},
	PSQL:                    {},
	InMemoryDB:              {},
	InMemoryDBV2:            {},
	Kafka:                   {},
	Logging:                 {},
	Mariadb:                 {},
	MariaDBV2:               {},
	Monitoring:              {},
	NFS:                     {},
	ObjectStorage:           {},
	VPN:                     {},
	PSQLV2:                  {},
}

// Validate reports every problem found in the file config as a single combined error,
// or nil when the config is valid (a nil receiver is treated as valid).
func (f *FileConfig) Validate() error {
	if f == nil {
		return nil
	}
	return errors.Join(
		f.validateProfilesAndEnvironments(),
		f.validateProducts(),
		f.Failover.Validate(),
	)
}

// validateProfilesAndEnvironments checks that environment and profile names are unique,
// and that currentProfile and each profile's environment point to entries that exist.
func (f *FileConfig) validateProfilesAndEnvironments() error {
	var problems []error

	// Check for duplicate environment names.
	envNames := make(map[string]struct{}, len(f.Environments))
	for _, env := range f.Environments {
		key := normalizeName(env.Name)
		if _, dup := envNames[key]; dup {
			problems = append(problems, fmt.Errorf("duplicate environment name %q", env.Name))
		}
		envNames[key] = struct{}{}
	}

	// Check for duplicate profile names and unknown environment references.
	profileNames := make(map[string]struct{}, len(f.Profiles))
	for _, profile := range f.Profiles {
		key := normalizeName(profile.Name)
		if _, dup := profileNames[key]; dup {
			problems = append(problems, fmt.Errorf("duplicate profile name %q", profile.Name))
		}
		profileNames[key] = struct{}{}

		if env := normalizeName(profile.Environment); env != "" {
			if _, envExists := envNames[env]; !envExists {
				problems = append(problems, fmt.Errorf("profile %q references unknown environment %q", profile.Name, profile.Environment))
			}
		}
	}

	// Check currentProfile refers to a defined profile.
	if current := normalizeName(f.CurrentProfile); current != "" {
		if _, ok := profileNames[current]; !ok {
			problems = append(problems, fmt.Errorf("currentProfile %q does not match any defined profile", f.CurrentProfile))
		}
	}

	return errors.Join(problems...)
}

// validateProducts checks, for every product in every environment, that it is a known
// product, has at least one endpoint, has no empty endpoint names, and is not declared
// more than once within the same environment.
func (f *FileConfig) validateProducts() error {
	var problems []error

	for _, env := range f.Environments {
		seenProduct := make(map[string]struct{}, len(env.Products))
		for _, product := range env.Products {
			key := normalizeName(product.Name)

			// Check for duplicate products in the environment.
			if _, dup := seenProduct[key]; dup {
				problems = append(problems, fmt.Errorf("duplicate product %q in environment %q", product.Name, env.Name))
			}
			seenProduct[key] = struct{}{}

			// Check the product name is recognized.
			if _, ok := knownProducts[key]; !ok {
				problems = append(problems, fmt.Errorf("unknown product %q in environment %q", product.Name, env.Name))
			}

			// Check the product declares at least one endpoint.
			if len(product.Endpoints) == 0 {
				problems = append(problems, fmt.Errorf("product %q in environment %q has no endpoints", product.Name, env.Name))
			}
			// Check for empty endpoint names.
			for i, ep := range product.Endpoints {
				if strings.TrimSpace(ep.Name) == "" {
					problems = append(problems, fmt.Errorf("product %q in environment %q has an endpoint with an empty name (position %d)", product.Name, env.Name, i+1))
				}
			}

			// A cloud product must not mix global (no location) and location-based endpoints.
			if key == Cloud {
				var hasGlobal, hasLocation bool
				for _, ep := range product.Endpoints {
					if strings.TrimSpace(ep.Location) == "" {
						hasGlobal = true
					} else {
						hasLocation = true
					}
				}
				if hasGlobal && hasLocation {
					problems = append(problems, fmt.Errorf("product %q in environment %q defines both global and location-based endpoints; only one of the two is allowed", product.Name, env.Name))
				}
			}
		}
	}

	return errors.Join(problems...)
}
