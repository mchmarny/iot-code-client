package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFloatProvider(t *testing.T) {

	pr, err := ParseProvider("temp|celsius|float|0:72.1|3s")
	assert.Nil(t, err)
	assert.NotNil(t, pr)

	rp := pr.GetParam()
	assert.NotNil(t, rp)
	assert.NotEmpty(t, rp.Label)
	assert.NotEmpty(t, rp.Unit)
	assert.NotEmpty(t, rp.Raw)

	assert.NotNil(t, rp.Frequency)

	assert.NotNil(t, rp.Template)
	assert.NotEmpty(t, rp.Template.Type)
	assert.NotNil(t, rp.Template.Min)
	assert.NotNil(t, rp.Template.Max)

}
