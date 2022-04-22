package utils

import (
	"testing"

	"bitbucket.org/level27/lvl/types"
)

// Test that there are no deserialization errors or similar from loading all domains
func TestDomainGetAll(t *testing.T) {
	client := makeTestClient()

	client.Domains(types.CommonGetParams{Limit: 1000000})
}