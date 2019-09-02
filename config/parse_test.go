package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	in := &Config{
		BaseImage: "foo.bar/alpine.tar.gz",
		InitCommands: []Command{{
			CommandLine:      "ls /root",
			Env:              nil,
			WorkingDirectory: "",
		}, {
			CommandLine:      "ls/bar",
			Env:              nil,
			WorkingDirectory: "",
		}},
		Copies: []Copies{{
			Source:      "/root/foo/bar",
			Destination: "/etc/systemd/system/",
		}},
	}

	bytes, err := yaml.Marshal(in)

	if err != nil {
		t.Fatal(err)
	}

	fd, err := os.OpenFile("conf.yaml", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := fd.Write(bytes); err != nil {
		t.Fatal(err)
	}
}
