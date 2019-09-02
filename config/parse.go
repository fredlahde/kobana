package config

import (
	e "github.com/fredlahde/kobana/errors"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func ParseConfig(path string) (*Config, error) {
	op := e.Op("config.ParseConfig")

	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, e.E(op, e.NotExists, err, e.C("config does not exist"), e.P("path", path))
	} else if os.IsPermission(err) {
		return nil, e.E(op, e.Permissions, err, e.C("can not read config, check permissions"), e.P("path", path))
	} else if stat.IsDir() {
		return nil, e.E(op, e.IsDir, err, e.C("path to config represents a directory"), e.P("path", path))
	}

	fd, err := os.Open(path)
	if err != nil {
		return nil, e.E(op, e.IO, err, e.C("can not read config"), e.P("path", path))
	}
	defer fd.Close()

	conf := &Config{}

	bytes, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, e.E(op, e.IO, err, e.C("failed to read config"), e.P("path", path))
	}

	err = yaml.Unmarshal(bytes, conf)
	if err != nil {
		return nil, e.E(op, e.Parse, err, e.C("failed to parse config"), e.P("path", path))
	}

	err = fd.Close()
	if err != nil {
		return nil, e.E(op, e.IO, err, e.C("failed to close config fd"), e.P("path", path))
	}

	return conf, nil
}
