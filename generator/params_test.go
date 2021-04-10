package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseParams(t *testing.T) {
	s := "path_prefix=/tmp,package_name=api"
	p, err := parseParams(s)
	assert.Nil(t, err)
	assert.Equal(t, p.PathPrefix, "/tmp")
}
