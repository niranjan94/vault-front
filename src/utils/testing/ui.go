package testing

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

import (
	"errors"
	"github.com/MrLYC/Golang-Which/src/command"
	"os/exec"
)

func findExecPath(commandName string) string {
	cmd := command.NewCommand(commandName)
	next := cmd.Find()
	for path := next(); len(path) > 0; path = next() {
		return path
	}
	return ""
}

func runNodeScript(script string, arguments ...string) (error) {
	nodePath := findExecPath("node")
	if nodePath == "" {
		return errors.New("could not find node executable")
	}
	cmd := exec.Command(nodePath, append([]string{script}, arguments...)...)
	cmd.Dir = os.Getenv("TEST_WORKING_DIR") + "/ui"

	stderr, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

	return cmd.Wait()
}

func RunUIE2E(t *testing.T) {
	assert.NoError(
		t,
		runNodeScript("node_modules/.bin/vue-cli-service", "test:e2e"),
	)
}
