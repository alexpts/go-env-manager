package config

import (
	"fmt"
)

type Config struct {
	GoPrivate string `yaml:"GOPRIVATE"`

	GoProxy   string `yaml:"GOPROXY"`
	GoNoProxy string `yaml:"GONOPROXY"`

	GoSumDB   string `yaml:"GOSUMDB"`
	GoNoSumDB string `yaml:"GONOSUMDB"`
}

func FromProcessEnv() Config {
	return Config{
		GoPrivate: GetGoEnv("GOPRIVATE"),
		GoProxy:   GetGoEnv("GOPROXY"),
		GoNoProxy: GetGoEnv("GONOPROXY"),
		GoSumDB:   GetGoEnv("GOSUMDB"),
		GoNoSumDB: GetGoEnv("GONOSUMDB"),
	}
}

func (c *Config) String() string {
	return fmt.Sprintf("GOPRIVATE: %s\nGOPROXY: %s\nGONOPROXY: %s\nGOSUMDB: %s\nGONOSUMDB: %s\n",
		c.GoPrivate,
		c.GoProxy,
		c.GoNoProxy,
		c.GoSumDB,
		c.GoNoSumDB,
	)
}
