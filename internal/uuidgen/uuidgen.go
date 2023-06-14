package uuidgen

import (
	uuid "github.com/gofrs/uuid/v5"
)

// Generates a UUID V5 for a resource.
func ResourceUuid() uuid.UUID {
	// Must is a helper that wraps a call to a function returning (UUID, error) and panics if the error is non-nil.
	uuidV4 := uuid.Must(uuid.NewV4())
	uuidV5 := uuid.NewV5(uuid.NewV5(uuid.NamespaceURL, "https://github.com/ionos-cloud/terraform-provider-ionoscloud"), uuidV4.String())

	return uuidV5
}
