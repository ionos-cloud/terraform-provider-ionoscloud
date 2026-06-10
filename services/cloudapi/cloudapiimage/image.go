package cloudapiimage

import (
	"context"
	"strings"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/cloudapilocation"
)

// Currently, this is not a service per se, but in the future, when the image service will be
// created, these functions can be included in the service. Right now, they are just utility
// functions in order to reuse the code.

// GetAllImages returns all images visible to the client's contract.
func GetAllImages(ctx context.Context, client *ionoscloud.APIClient) ([]ionoscloud.Image, error) {
	images, apiResponse, err := client.ImagesApi.ImagesGet(ctx).Depth(1).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, err
	}
	if images.Items == nil {
		return nil, nil
	}
	return *images.Items, nil
}

// GetImageAlias returns the alias matching want carried by any of the given images in one
// of the given locations, in canonical casing, or "" when none matches. Aliases are
// matched from the images rather than the location's alias catalog, which restricted
// contracts may not be able to read. Image type is deliberately ignored: e.g. "*:*_iso"
// aliases live on CDROM images, and the alias string is what gets sent to the API.
func GetImageAlias(want string, images []ionoscloud.Image, locations []string) string {
	if want == "" {
		return ""
	}
	for _, img := range images {
		if img.Properties == nil || img.Properties.ImageAliases == nil ||
			img.Properties.Location == nil || !cloudapilocation.LocationInSet(locations, *img.Properties.Location) {
			continue
		}
		if alias := MatchImageAlias(*img.Properties.ImageAliases, want); alias != "" {
			return alias
		}
	}
	return ""
}

// MatchImageAlias returns the alias from aliases that case-insensitively matches want,
// preserving the alias' canonical casing, or "" when none matches.
func MatchImageAlias(aliases []string, want string) string {
	if want == "" {
		return ""
	}
	for _, alias := range aliases {
		if alias != "" && strings.EqualFold(alias, want) {
			return alias
		}
	}
	return ""
}
