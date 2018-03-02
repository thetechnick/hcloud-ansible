package version

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	VersionTag = "1.2.3"
	Branch = "test"
	BuildDate = "1512490612"

	assert := assert.New(t)
	v := Get()

	assert.Equal(VersionTag, v.Version)
	assert.Equal(Branch, v.Branch)
	assert.Equal(BuildDate, fmt.Sprint(v.BuildDate.Unix()))
}
