package config

type LogConfig struct {
	Type   string            `json:"log-driver,omitempty"`
	Config map[string]string `json:"log-opts,omitempty"`
}

// Config defines the configuration of a docker daemon.
// It includes json tags to deserialize configuration from a file
// using the same names that the flags in the command line uses.
type Config struct {
	// Fields below here are platform specific.
	LogConfig
	LogLevel string
	Debug    bool
	RawLogs  bool `json:"raw-logs,omitempty"`
	Rootless bool `json:"rootless,omitempty"`
}

// IsRootless returns conf.Rootless
func (conf *Config) IsRootless() bool {
	return conf.Rootless
}
func New() *Config {
	config := Config{}
	config.LogConfig.Config = make(map[string]string)
	return &config
}
