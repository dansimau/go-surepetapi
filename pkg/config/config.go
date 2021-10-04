package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dansimau/go-surepetapi"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	API surepetapi.Config `yaml:"api"`
}

// FromFile loads app config from a YAML file.
func FromFile(file string) (*Config, error) {
	yamlBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config := Config{}
	if err := yaml.Unmarshal(yamlBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// searchPaths returns a list of paths from basePath upwards to the root ("/).
func searchPaths(basePath string) (paths []string) {
	root := basePath

	for root != "/" {
		paths = append(paths, root)
		root = filepath.Dir(root)
	}

	paths = append(paths, "/")

	return paths
}

func searchParentsForPath(filename string) (path string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for _, path := range searchPaths(wd) {
		f := filepath.Join(path, filename)
		if fileExists(f) {
			return f, nil
		}
	}

	return "", nil
}

func SearchAndLoadConfig(filename string) (cfg *Config, path string, err error) {
	configPath, err := searchParentsForPath(filename)
	if err != nil {
		return nil, "", err
	}

	if configPath == "" {
		return nil, "", errors.New("cannot find config")
	}

	cfg, err = FromFile(configPath)
	if err != nil {
		return nil, "", err
	}

	return cfg, configPath, nil
}
