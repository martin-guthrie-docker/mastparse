package mastparse

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
)

// initializer struct for ConfigClass
// - items from the extern yaml config file, or the environment vars
type ConfigClassCfg struct {
	Log                  *logrus.Logger

	// add fields from config file here
	MastPath             string

	// environment variable data fields here
}


type ConfigClass struct {
	ConfigClassCfg       // this is an embedded type

	// add other fields here
}

// constructor for ConfigClass
func NewConfigClass(v *viper.Viper, cfg ConfigClassCfg) (*ConfigClass, error) {

	// if no logger, create a null logger
	if cfg.Log == nil {
		cfg.Log = logrus.New()
		cfg.Log.Out = ioutil.Discard
	}

	t := new(ConfigClass)
	t.Log = cfg.Log

	// from viper, if there is any data in input config file
	err := v.Unmarshal(&t.ConfigClassCfg)
	if err != nil {
		panic(err)
	}

	// data from the environment variables
	// was not able to use Unmarshal.... see https://github.com/spf13/viper/issues/188
	//if v.Get("TBD") != nil {
	//	t.ConfigClass.TBD = v.Get("one").(string)
	//} else {
	//	t.ConfigClass.TBD = "UNKNOWN"
	//}

	return t, nil
}

func (t *ConfigClass) Dump(toConsole bool) error {
	// dump anything as needed
	t.Log.Infof("ConfigClassCfg.MastPath: %s", t.ConfigClassCfg.MastPath)
	if toConsole {
		fmt.Printf("ConfigClassCfg.MastPath: %s\n", t.ConfigClassCfg.MastPath)
	}

	return nil
}
