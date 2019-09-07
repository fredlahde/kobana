package config

type Command struct {
	CommandLine      string   `yaml:"cmd"`
	Env              []string `yaml:"env"`
	WorkingDirectory string   `yaml:"workingDir"`
}

func (cmd *Command) Equal(other *Command) bool {
	if cmd.CommandLine != other.CommandLine {
		return false
	}

	if cmd.WorkingDirectory != other.WorkingDirectory {
		return false
	}

	for i, e := range cmd.Env {
		oe := other.Env[i]
		if e != oe {
			return false
		}
	}

	return true
}

type Mapping struct {
	Source      string `yaml:"src"`
	Destination string `yaml:"dst"`
}

func (m *Mapping) Equal(other *Mapping) bool {
	return m.Source == other.Source && m.Destination == other.Destination
}

type Config struct {
	BaseImage    string    `yaml:"base"`     // linux base image to use
	InitCommands []Command `yaml:"commands"` // commands to run to setup the chroot
	Mappings     []Mapping `yaml:"mappings"` // readonly copies from host to chroot
	Env          []string  `yaml:"env"`
	EnvFile      []string  `yaml:"env_file"`
}

func (c *Config) Equal(other *Config) bool {
	if c.BaseImage != other.BaseImage {
		return false
	}

	for i, e := range c.Env {
		other := other.Env[i]
		if e != other {
			return false
		}
	}

	for i, e := range c.EnvFile {
		other := other.EnvFile[i]
		if e != other {
			return false
		}
	}

	for i, mapping := range c.Mappings {
		other := other.Mappings[i]

		if !mapping.Equal(&other) {
			return false
		}
	}

	for i, cmd := range c.InitCommands {
		oCmd := other.InitCommands[i]

		if !cmd.Equal(&oCmd) {
			return false
		}
	}

	return true
}
