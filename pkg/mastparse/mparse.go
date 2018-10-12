package mastparse

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)


const mastInventoryFile string = "/inventory/1.hosts"
const ansibleVarsRegex string = `ansible_user=(?P<user>\w*)\s*ansible_host=(?P<ip>\d+\.\d+\.\d+\.\d+)`

// initializer struct for ExAClass
type MparseClassCfg struct {
	Log                  *logrus.Logger
	Name 	 		     string  // name of deployment
	MastPath             string  // path of mast information, typically ~/.mast
}

type creds struct {
	name    string
	user    string
	ip      string
}


type MparseClass struct {
	MparseClassCfg          // this is an embedded type

	// internal state vars
	openSucceeded        bool
	pathToMastInventory  string

	// things needed to build out the commands
	ansibleRegex         *regexp.Regexp
	docker_ucp_lb        string
	sshHosts             []creds
	hosts                []string

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
	t.ansibleRegex, _ = regexp.Compile(ansibleVarsRegex)

	t.hosts = []string{"linux-ucp-manager-primary", "linux-dtr-worker-primary"}

	return t, nil
}

func (t *MparseClass) get_creds_fields(dst *[]creds, body string) error {
	for _, line := range strings.Split(strings.TrimSuffix(body, "\n"), "\n") {
		t.Log.Debug("parse: ", line)
		matches := t.ansibleRegex.FindStringSubmatch(line)
		if len(matches) != 3 {
			t.Log.Errorf("Unexpected matches")
			return errors.New("Unexpected matches")
		}
		t.Log.Debug("len(matches): ", len(matches))
		t.Log.Debug("matches[0]: ", matches[0])
		t.Log.Info("ucp_manager_user: ", matches[1])
		t.Log.Info("ucp_manager_ip: ", matches[2])
		var c creds
		c.user = matches[1]
		c.ip = matches[2]
		*dst = append(*dst, c)
	}
	return nil
}

func (t *MparseClass) get_body(field string) (string, error) {
	cfg, err := ini.LoadSources(ini.LoadOptions{
		UnparseableSections: []string{field},
	}, t.pathToMastInventory)
	if err != nil {
		t.Log.Errorf("Fail to ini.LoadSources: ", err)
		return "", err
	}

	body := cfg.Section(field).Body()
	return body, nil
}

func (t *MparseClass) get_hosts() error {
	// the ansible ini field looks like this,
	//
	// [linux-ucp-manager-primary]
	// i-0dc78492b42702fa5 ansible_user=docker ansible_host=54.202.38.187
	//
	// which is not parsable via the ini library, so bring in as a blob
	// and use regex to get the parts

	for _, host := range t.hosts {
		body, err := t.get_body(host)
		if err != nil {
			t.Log.Errorf("Fail to ini.get_body for host: %s, %v", host, err)
			return err
		}
		t.get_creds_fields(&t.sshHosts, body)
	}
	t.Log.Debug("len(t.sshHosts): ", len(t.sshHosts))

	return nil
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
	t.Log.Info("docker_ucp_lb: ", t.docker_ucp_lb)

	err = t.get_hosts()
	if err != nil {
		t.Log.Errorf("Fail get_ucp_manager_prime_fields: %v", err)
		return err
	}

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

