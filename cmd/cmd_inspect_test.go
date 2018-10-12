package cmd_test

import (
	"testing"

	"github.com/martin-guthrie-docker/mastparse/cmd"
)

func TestAction1Func(t *testing.T) {

	_, err := cmd.ExecuteCommand("deployment", "string1")

	if err != nil {
		t.Errorf("deployment failed with err: %s", err.Error())
	}
}
