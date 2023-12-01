package ptr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	k := Knife{
		Name: "k0",
	}

	d := NewDecider(&k, "knife")

	assert.Equal(t, false, d.Switch.GetState())

	d.Switch.Toggle()

	assert.Equal(t, true, d.Switch.GetState())

}
