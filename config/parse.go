package config

import (
	e "github.com/fredlahde/kobana/errors"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var globalConfig *Config

func Get() *Config {
	return globalConfig
}

func ParseConfig(path string) error {
	op := e.Op("globalConfig.ParseConfig")

	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return e.E(op, e.NotExists, err, e.C("globalConfig does not exist"), e.P("path", path))
	} else if os.IsPermission(err) {
		return e.E(op, e.Permissions, err, e.C("can not read globalConfig, check permissions"), e.P("path", path))
	} else if stat.IsDir() {
		return e.E(op, e.IsDir, err, e.C("path to globalConfig represents a directory"), e.P("path", path))
	}

	fd, err := os.Open(path)
	if err != nil {
		return  e.E(op, e.IO, err, e.C("can not read globalConfig"), e.P("path", path))
	}
	defer fd.Close()

	conf := &Config{}

	bytes, err := ioutil.ReadAll(fd)
	if err != nil {
		return e.E(op, e.IO, err, e.C("failed to read globalConfig"), e.P("path", path))
	}

	err = yaml.Unmarshal(bytes, conf)
	if err != nil {
		return e.E(op, e.Parse, err, e.C("failed to parse globalConfig"), e.P("path", path))
	}

	err = fd.Close()
	if err != nil {
		return e.E(op, e.IO, err, e.C("failed to close globalConfig fd"), e.P("path", path))
	}

	globalConfig = conf
	return nil
}
