package mastparse

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)


const mastInventoryFile string = "/inventory/1.hosts"
const ansibleVarsRegex string = `^(?P<id>.*)\s*ansible_user=(?P<user>\w*)\s*ansible_host=(?P<ip>\d+\.\d+\.\d+\.\d+)`

// initializer struct for ExAClass
type MparseClassCfg struct {
	Log                  *logrus.Logger
	Name 	 		     string  // name of deployment
	MastPath             string  // path of mast information, typically ~/.mast
}

type creds struct {
	field   string
	id      string
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
	iniField             []string

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

	t.iniField = []string{"linux-ucp-manager-primary", "linux-dtr-worker-primary", "linux-ucp-manager-replicas",
	"linux-dtr-worker-replicas", "linux-workers", "windows-workers", "linux-databases", "linux-build-servers",
	"windows-databases", "windows-build-servers"}

	return t, nil
}

func (t *MparseClass) get_creds_fields(dst *[]creds, body string, field string) error {
	t.Log.Info("parsing field: ", field)
	for _, line := range strings.Split(strings.TrimSuffix(body, "\n"), "\n") {
		if len(line) == 0 {
			continue
		}
		t.Log.Debug("line: ", line)
		matches := t.ansibleRegex.FindStringSubmatch(line)
		if len(matches) == 0 {
			t.Log.Debugf("No matches for line: %s", line)
			continue
		}
		if len(matches) != 4 {
			t.Log.Errorf("Unexpected matches")
			continue
		}
		t.Log.Debug("len(matches): ", len(matches))
		t.Log.Debug("matches[0]: ", matches[0])
		t.Log.Info("id  : ", matches[1])
		t.Log.Info("user: ", matches[2])
		t.Log.Info("ip  : ", matches[3])
		var c creds
		c.field = field
		c.id = matches[1]
		c.user = matches[2]
		c.ip = matches[3]
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

	for _, field := range t.iniField {
		body, err := t.get_body(field)
		if err != nil {
			t.Log.Errorf("Fail to ini.get_body for host: %s, %v", field, err)
			return err
		}
		if len(body) == 0 {
			continue
		}
		t.get_creds_fields(&t.sshHosts, body, field)
	}
	t.Log.Debug("len(t.sshHosts): ", len(t.sshHosts))

	return nil
}

func (t *MparseClass) make_ssh_for_hosts() error {
	for _, cred := range t.sshHosts {
		fmt.Printf("%25s : %20s : ssh -i ~./mast/id_rsa %s@%s\n", cred.field, cred.id, cred.user, cred.ip)
	}

	return nil
}

func (t *MparseClass) ReadMastInventory() error {
	t.Log.Info("Start: path ", t.pathToMastInventory)

	if !t.openSucceeded {
		t.Log.Error("Open not succeeded")
		return errors.New("Open not succeeded")
	}

	// Easy stuff, parse the ini file for easy fields we care about
	// load the file...
	cfg, err := ini.Load(t.pathToMastInventory)
	if err != nil {
		t.Log.Errorf("Fail to read file: %v", err)
		return err
	}

	// get section and variable
	t.docker_ucp_lb = cfg.Section("all:vars").Key("docker_ucp_lb").String()
	t.Log.Info("docker_ucp_lb: ", t.docker_ucp_lb)

	// Hard stuff, ansible ini file has lines that don't parse as ini
	// find all the hosts for ssh creds
	err = t.get_hosts()
	if err != nil {
		t.Log.Errorf("Fail get_hosts: %v", err)
		return err
	}
	err = t.make_ssh_for_hosts()
	if err != nil {
		t.Log.Errorf("Fail make_ssh_cmds: %v", err)
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

