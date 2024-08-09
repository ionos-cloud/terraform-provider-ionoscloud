package uuidgen

import (
	uuid "github.com/gofrs/uuid/v5"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// Generates a UUID V5 for a resource.
func ResourceUuid() uuid.UUID {
	// Must is a helper that wraps a call to a function returning (UUID, error) and panics if the error is non-nil.
	uuidV4 := uuid.Must(uuid.NewV4())
	uuidV5 := uuid.NewV5(uuid.NewV5(uuid.NamespaceURL, constant.RepoURL), uuidV4.String())

	return uuidV5
}

func GenerateUuidFromName(name string) string {
	return uuid.NewV5(uuid.NewV5(uuid.NamespaceURL, constant.RepoURL), name).String()
}
