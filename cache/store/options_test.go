package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptionsCostValue(t *testing.T) {
	// Given: an Options struct with a specific Cost value
	options := &Options{
		Cost: 7,
	}

	// When: retrieving the Cost value
	actualCost := options.Cost

	// Then: the Cost value should match the expected value
	assert.Equal(t, int64(7), actualCost, "The Cost value should be 7")
}

func TestOptionsExpirationValue(t *testing.T) {
	// Given: an Options struct with a specific Expiration value
	options := &Options{
		Expiration: 25 * time.Second,
	}

	// When: retrieving the Expiration value
	actualExpiration := options.Expiration

	// Then: the Expiration value should match the expected value
	assert.Equal(t, 25*time.Second, actualExpiration, "The Expiration value should be 25 seconds")
}

func TestOptionsTagsValue(t *testing.T) {
	// Given: an Options struct with specific Tags
	options := &Options{
		Tags: []string{"tag1", "tag2", "tag3"},
	}

	// When: retrieving the Tags
	actualTags := options.Tags

	// Then: the Tags should match the expected values
	assert.Equal(t, []string{"tag1", "tag2", "tag3"}, actualTags, "The Tags should be ['tag1', 'tag2', 'tag3']")
}

func TestApplyOptionsWithDefault(t *testing.T) {
	// Given: default Options with a specific Expiration value
	defaultOptions := &Options{
		Expiration: 25 * time.Second,
	}

	// When: applying additional options to the default options
	options := ApplyOptionsWithDefault(defaultOptions, WithCost(7))

	// Then: the resulting Options should have the expected Cost and Expiration values
	assert.Equal(t, int64(7), options.Cost, "The Cost value should be 7")
	assert.Equal(t, 25*time.Second, options.Expiration, "The Expiration value should be 25 seconds")
}
