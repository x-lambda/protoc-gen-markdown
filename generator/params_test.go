package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseParams(t *testing.T) {
	s := "prefix=/tmp,package=demo.v0,filename=out.md"
	p, err := ParseParams(s)
	assert.Nil(t, err)

	assert.Equal(t, p.Prefix, "/tmp")
	assert.Equal(t, p.Package, "demo.v0")
	assert.Equal(t, p.Filename, "out.md")
}
