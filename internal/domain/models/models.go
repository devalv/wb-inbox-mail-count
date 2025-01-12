package models

type ServerConfig struct {
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	Debug   bool           `yaml:"debug"`
	Servers []ServerConfig `yaml:"servers"`
}

// TODO: WaybarOutput
