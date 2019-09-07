package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	in := &Config{
		BaseImage: "http://dl-cdn.alpinelinux.org/alpine/v3.10/releases/x86_64/alpine-minirootfs-3.10.2-x86_64.tar.gz",
		EnvFile:   []string{".env"},
		Env:       []string{"FOO=BAR"},
		InitCommands: []Command{{
			CommandLine:      "ls /root",
			Env:              []string{"GOOS=foo"},
			WorkingDirectory: "",
		}, {
			CommandLine:      "ls/bar",
			Env:              nil,
			WorkingDirectory: "",
		}},
		Mappings: []Mapping{{
			Source:      "/root/foo/bar",
			Destination: "/etc/systemd/system/",
		}, {
			Source:      "$HOME/foo.txt",
			Destination: "foo.txt",
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

	var out Config
	yaml.Unmarshal(bytes, &out)
	if !in.Equal(&out) {
		t.Fatal("structs do not match")
	}
}
