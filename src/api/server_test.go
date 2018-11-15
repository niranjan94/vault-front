package api

import (
	"context"
	"github.com/GeertJohan/go.rice"
	"github.com/niranjan94/vault-front/src/cmd"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	cmd.LoadConfigForTest()
}

func TestStartApiServer(t *testing.T) {
	e := StartApiServer(
		rice.MustFindBox("../../ui/dist").HTTPBox(),
		false,
	)
	assert.NotEmpty(t, e)
	// testingUtils.RunUIE2E(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}