package vault

import (
	"testing"
)

func TestGetSelfEntityId(t *testing.T) {
	client := GetManagerClient()
	GetSelfEntityId(client)
}