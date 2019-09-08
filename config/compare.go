package config

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

func (m *Mapping) Equal(other *Mapping) bool {
	return m.Source == other.Source && m.Destination == other.Destination
}

func (c *Config) Equal(other *Config) bool {
	if c.BaseImage != other.BaseImage {
		return false
	}

	if c.BaseDir != other.BaseDir {
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
