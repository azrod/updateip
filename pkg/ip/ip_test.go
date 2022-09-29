package ip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetMyExternalIP(t *testing.T) {
	ip, err := GetMyExternalIP()
	assert.Nil(t, err)
	assert.NotEmpty(t, ip.String())
}
