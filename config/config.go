package config

type Config struct {
	Port        int
	AllowedExts map[string]bool
}

func NewConfig(port int, exts map[string]bool) *Config {
	return &Config{
		Port:        port,
		AllowedExts: exts,
	}
}
