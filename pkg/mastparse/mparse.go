package mastparse

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
)


const mastInventoryFile string = "/inventory/1.hosts"

// initializer struct for ExAClass
type MparseClassCfg struct {
	Log                  *logrus.Logger
	Name 	 		     string  // name of deployment
	MastPath             string  // path of mast information, typically ~/.mast
}

type MparseClass struct {
	MparseClassCfg          // this is an embedded type

	// internal state vars
	openSucceeded        bool
	pathToMastInventory  string

	// things needed to build out the commands
	docker_ucp_lb        string
	ucp_manager_ip       string
	ucp_manager_user     string
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

	// set internal states
	t.openSucceeded = false

	return t, nil
}

func (t *MparseClass) ReadMastInventory() error {
	t.Log.Info("Start: path ", t.pathToMastInventory)

	if !t.openSucceeded {
		t.Log.Error("Open not succeeded")
		return errors.New("Open not succeeded")
	}

	cfg, err := ini.Load(t.pathToMastInventory)
	if err != nil {
		t.Log.Errorf("Fail to read file: %v", err)
		return err
	}

	t.docker_ucp_lb = cfg.Section("all:vars").Key("docker_ucp_lb").String()
	t.Log.Info("docker_ucp_lb:", t.docker_ucp_lb)

	keys := cfg.Section("linux-ucp-manager-primary").KeyStrings()
	t.Log.Debug(keys[0])

	value := cfg.Section("linux-ucp-manager-primary").Key(keys[0]).Value()
	t.Log.Debug(value)
	t.Log.Info("test")

	ansible_info := cfg.Section("linux-ucp-manager-primary").Key(keys[0]).String()
	t.Log.Info("ucp_manager_ip:", ansible_info)

	return nil
}


func (t *MparseClass) Open() error {
	t.Log.Info("Start: name ", t.Name, ", path ", t.MastPath)

	// check is mast path exist
	if _, err := os.Stat(t.MastPath); os.IsNotExist(err) {
		// path/to/whatever does not exist
		t.Log.Error(t.MastPath, " does not exist")
		return err
	}

	// check is deployment Name path exist
	if _, err := os.Stat(t.MastPath + "/" + t.Name); os.IsNotExist(err) {
		// path/to/whatever does not exist
		t.Log.Error(t.MastPath, " does not exist")
		return err
	}
	t.pathToMastInventory = t.MastPath + "/" + t.Name + mastInventoryFile

	// check is named inventory path exist
	if _, err := os.Stat(t.pathToMastInventory); os.IsNotExist(err) {
		// path/to/whatever does not exist
		t.Log.Error(t.MastPath, " does not exist")
		return err
	}

	t.Log.Info(t.pathToMastInventory + " is present")
	t.openSucceeded = true

	return nil
}

func (t *MparseClass) OpenSucceeded() bool {
	return t.openSucceeded
}

func (t *MparseClass) Close() {
	t.Log.Info("Start")
}

