package mastparse

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

// initializer struct for ExAClass
type MparseClassCfg struct {
	Log                  *logrus.Logger
	Name 	 		     string  // name of deployment
	MastPath             string  // path of mast information, typically ~/.mast
}

type MparseClass struct {
	MparseClassCfg          // this is an embedded type
}

// constructor for ExAClass
func NewMparseClass(cfg MparseClassCfg) (*MparseClass, error) {

	// if no logger, create a null logger
	if cfg.Log == nil {
		cfg.Log = logrus.New()
		cfg.Log.Out = ioutil.Discard
	}

	t := new(MparseClass)

	// transfer config settings to this class
	// TODO: how to do this in one line... like a memcopy
	t.Log = cfg.Log
	t.Name = cfg.Name
	t.MastPath = cfg.MastPath

	return t, nil
}

func (t *MparseClass) Open() error {
	t.Log.Info("Start: name ", t.Name, ", path ", t.MastPath)
	return nil
}

func (t *MparseClass) Close() {
	t.Log.Info("Start")
}