package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSHA512(t *testing.T) {
	hash := SHA512("Always code as if the guy who ends up maintaining your code will be a violent psychopath who knows where you live")
	assert.Equal(t, "c9616d3f8e43441feca673663c40e71ba6b9593a33ca331a224881f6952fd619ebf16180f5474cf01509fd203c24294a3ae1706b287a360100cb422b9c406fa9", hash)
}