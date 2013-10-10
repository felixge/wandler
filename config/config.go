package config

import (
	"io/ioutil"
	"launchpad.net/goyaml"
)

func ReadFile(path string, config interface{}) (error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return goyaml.Unmarshal(data, config)
}
