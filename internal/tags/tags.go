package tags

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	objstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
)

// KeyValueTags is a map of key-value tags.
type KeyValueTags map[string]string

// New creates a new KeyValueTags from a list of objstorage.Tag.
func New(tags []objstorage.Tag) KeyValueTags {
	result := make(KeyValueTags)

	for _, tag := range tags {
		result[tag.Key] = tag.Value
	}

	return result
}

// ToMap converts KeyValueTags to a Terraform map.
func (t KeyValueTags) ToMap(ctx context.Context) (types.Map, error) {
	if len(t) == 0 {
		return types.MapNull(types.StringType), nil
	}

	tfResult, diagErr := types.MapValueFrom(ctx, types.StringType, t)
	if diagErr != nil {
		return types.Map{}, fmt.Errorf("failed to convert KeyValueTags to types.Map: %v", diagErr)
	}

	return tfResult, nil
}

// NewFromMap creates a new KeyValueTags from a Terraform map.
func NewFromMap(m types.Map) KeyValueTags {
	result := make(KeyValueTags)

	for k, v := range m.Elements() {
		result[k] = v.(types.String).ValueString()
	}

	return result
}

// Merge merges two KeyValueTags.
func (t KeyValueTags) Merge(other KeyValueTags) KeyValueTags {
	result := make(KeyValueTags)

	for k, v := range t {
		result[k] = v
	}

	for k, v := range other {
		result[k] = v
	}

	return result
}

// Ignore removes tags from KeyValueTags.
func (t KeyValueTags) Ignore(ignoreTags KeyValueTags) KeyValueTags {
	result := make(KeyValueTags)

	for k, v := range t {
		if _, ok := ignoreTags[k]; ok {
			continue
		}

		result[k] = v
	}

	return result
}

// ToList converts KeyValueTags to a list of objstorage.Tag.
func (t KeyValueTags) ToList() []objstorage.Tag {
	tags := make([]objstorage.Tag, 0, len(t))
	for key, value := range t {
		tags = append(tags, objstorage.Tag{Key: key, Value: value})
	}

	return tags
}

// ToListPointer converts KeyValueTags to a pointer to a list of objstorage.Tag.
func (t KeyValueTags) ToListPointer() *[]objstorage.Tag {
	tags := t.ToList()
	return &tags
}
