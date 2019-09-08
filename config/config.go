package config

type Command struct {
	CommandLine      string   `yaml:"cmd"`
	Env              []string `yaml:"env"`
	WorkingDirectory string   `yaml:"workingDir"`
}

type Mapping struct {
	Source      string `yaml:"src"`
	Destination string `yaml:"dst"`
}

type Config struct {
	BaseImage    string    `yaml:"base"`     // linux base image to use
	BaseDir      string    `yaml:"base_dir"` // base directory for chroot mounts
	InitCommands []Command `yaml:"commands"` // commands to run to setup the chroot
	Mappings     []Mapping `yaml:"mappings"` // readonly copies from host to chroot
	Env          []string  `yaml:"env"`
	EnvFile      []string  `yaml:"env_file"`
}
