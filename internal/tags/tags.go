package tags

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	s3 "github.com/ionos-cloud/sdk-go-s3"
)

// KeyValueTags is a map of key-value tags.
type KeyValueTags map[string]string

// New creates a new KeyValueTags from a list of s3.Tag.
func New(tags []s3.Tag) KeyValueTags {
	result := make(KeyValueTags)

	for _, tag := range tags {
		result[*tag.Key] = *tag.Value
	}

	return result
}

// NewFromTFMap creates a new KeyValueTags from a Terraform map.
func NewFromTFMap(m types.Map) KeyValueTags {
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

// ToList converts KeyValueTags to a list of s3.Tag.
func (t KeyValueTags) ToList() []s3.Tag {
	tags := make([]s3.Tag, 0, len(t))
	for key, value := range t {
		tags = append(tags, s3.Tag{Key: s3.PtrString(key), Value: s3.PtrString(value)})
	}

	return tags
}

// ToListPointer converts KeyValueTags to a pointer to a list of s3.Tag.
func (t KeyValueTags) ToListPointer() *[]s3.Tag {
	tags := t.ToList()
	return &tags
}
