package kafka

import (
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
)

func EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewUserCredentialsEphemeral,
	}
}
