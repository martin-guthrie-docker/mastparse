package cmd_test

import (
	"testing"

	"github.com/martin-guthrie-docker/mastparse/cmd"
)

func TestInspect(t *testing.T) {

	_, err := cmd.ExecuteCommand("inspect", "test")

	if err != nil {
		t.Errorf("deployment failed with err: %s", err.Error())
	}
}
