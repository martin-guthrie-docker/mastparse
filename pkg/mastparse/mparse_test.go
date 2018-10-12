package mastparse

import (
	"github.com/martin-guthrie-docker/mastparse/pkg/log"

	"testing"
)

const mastPath string = "./test_assets"

func TestNewMparseClass(t *testing.T) {
	log.Term.Info("Start ")

	mp, err := NewMparseClass(
		MparseClassCfg{
			Log:  log.Term,
			Name: "test",
			MastPath: mastPath,
		})

	if err != nil {
		t.Errorf("NewMparseClass failed: %s", err.Error())
	}

	if mp.OpenSucceeded() {
		t.Error("openSucceeded before Open() called")
	}
}

