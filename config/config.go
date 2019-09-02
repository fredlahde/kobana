package config

type Command struct {
	CommandLine      string   `yaml:"cmd"`
	Env              []string `yaml:"env"`
	WorkingDirectory string   `yaml:"workingDir"`
}

type Copies struct {
	Source      string `yaml:"src"`
	Destination string `yaml:"dst"`
}

type Config struct {
	BaseImage    string    `yaml:"base"`     // linux base image to use
	InitCommands []Command `yaml:"commands"` // commands to run to setup the chroot
	Copies       []Copies  `yaml:"copies"`   // readonly copies from host to chroot
}
