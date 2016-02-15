// Tests
package wrapper

import (
	"crypto/md5"
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"
)

var (
	cmd_rmrf = []string{
		"rm -r -f /",
		"rm -rf /",
		"rm -f -r /",
		"rm -rf /etc",
		"rm --recursive --force /",
	}

	cmd_ls = []string{
		"ls -r -f /",
		"ls -rf /",
		"ls -f -r /",
		"ls -rf /etc",
	}
)

func TestGoodCmd(t *testing.T) {
	f := Forbidden{}

	err := ReadForbidden("forbidden.yml", &f)
	if err != nil {
		t.Error(err)
		return
	}

	for _, value := range cmd_ls {
		_, ok, err := Filter(value, f)
		if err != nil {
			t.Error(err)
		}
		if !ok {
			t.Error("Cached", value)
		}
	}
}

func TestBadCmd(t *testing.T) {
	f := Forbidden{}

	err := ReadForbidden("forbidden.yml", &f)
	if err != nil {
		t.Error(err)
		return
	}

	for _, value := range cmd_rmrf {
		_, ok, err := Filter(value, f)
		if err != nil {
			t.Error(err)
		}
		if ok {
			t.Error("Not cached", value)
		}
	}
}

func TestReadForbidden(t *testing.T) {

	bs, err := ioutil.ReadFile("forbidden.yml")
	if err != nil {
		t.Error("Cannot open control yaml file")
	}

	m := Forbidden{}

	ReadForbidden("forbidden.yml", &m)

	st, err := yaml.Marshal(&m)
	if err != nil {
		t.Error("Unable to marshal")
	}

	if md5.Sum(st) != md5.Sum(bs) {
		t.Error("Input and output of yaml not equal")
	}
}
