package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type (
	TLS struct {
		Enable           bool   `yaml:"enable"`
		Listen           string `yaml:"listen"`
		CertificateChain string `yaml:"certificate_chain"`
		PrivateKey       string `yaml:"private_key"`
	}

	Listener struct {
		Listen string `yaml:"listen"`
		TLS    TLS    `yaml:"tls"`
	}

	Device struct {
		Name string `yaml:"name"`
		Baud int    `yaml:"baud"`
	}
)

var grid650ArraySerialConfig string

type Grid650ArraySerialConfig struct {
	HTTP   Listener `yaml:"http"`
	Device Device   `yaml:"device"`
}

func LoadGrid650ArraySerialConfig(file string) (config *Grid650ArraySerialConfig, err error) {
	// read config file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}
