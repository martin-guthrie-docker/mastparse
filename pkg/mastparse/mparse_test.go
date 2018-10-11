package mastparse

import (
	"github.com/sirupsen/logrus"
	"testing"
)

var log *logrus.Logger = logrus.New()

func TestFooClass_New(t *testing.T) {

	var str1 = "string1"
	var str2 = "string2"

	foo, err := NewFooClass(
		str1,
		 FooClassCfg{
			Name: str2,
			Log: log,
		 })
	if err != nil {
		t.Errorf("NewFooClass failed with err: %s", err.Error())
	}

	if foo.Host != str1 {
		t.Errorf("NewFooClass foo.Host failed %s != %s", foo.Host, str1)
	}

	if foo.Name != str2 {
		t.Errorf("NewFooClass foo.Name failed %s != %s", foo.Name, str2)
	}

	if foo.open {
		t.Errorf("NewFooClass foo.open failed %t", foo.open)
	}
}

func TestFooClass_Open(t *testing.T) {

	var str1 = "string1"
	var str2 = "string2"

	foo, err := NewFooClass(
		str1,
		FooClassCfg{
			Name: str2,
			Log: log,
		})
	if err != nil {
		t.Errorf("NewFooClass failed with err: %s", err.Error())
	}

	foo.Open()

	if !foo.open {
		t.Errorf("NewFooClass foo.open failed %t", foo.open)
	}
}

func TestFooClass_Close(t *testing.T) {

	var str1 = "string1"
	var str2 = "string2"

	foo, err := NewFooClass(
		str1,
		FooClassCfg{
			Name: str2,
			Log: log,
		})
	if err != nil {
		t.Errorf("NewFooClass failed with err: %s", err.Error())
	}

	foo.Open()
	foo.Close()

	if foo.open {
		t.Errorf("NewFooClass foo.open failed %t", foo.open)
	}
}