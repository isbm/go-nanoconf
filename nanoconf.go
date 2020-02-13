package nanoconf

import (
	"errors"
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Inspector struct {
	subset *map[string]interface{}
}

func NewInspector(tree *map[string]interface{}) *Inspector {
	ins := new(Inspector)
	ins.subset = tree
	return ins
}

// Return the entire tree raw
func (ins *Inspector) Raw() *map[string]interface{} {
	return ins.subset
}

// String returns a string type of a config value.
func (ins *Inspector) String(key string, overlay string) string {
	if overlay != "" {
		return overlay
	}
	return fmt.Sprintf("%s", (*ins.subset)[key])
}

// Int returns an integer type of a config value.
func (ins *Inspector) Int(key string, overlay string) (int, error) {
	var v string
	var err error
	var data int
	if overlay != "" {
		v = overlay
	} else {
		if (*ins.subset)[key] != nil {
			v = fmt.Sprintf("%d", (*ins.subset)[key])
		} else {
			err = errors.New("Value not found")
		}
	}
	if err == nil {
		data, err = strconv.Atoi(v)
	}

	return data, err
}

// DefaultInt is a wrapper around Int method, allowing return default value,
// in case nothing has been found. Overlay value is an empty string and it
// doesn't mean 0, but nil.
func (ins *Inspector) DefaultInt(val string, overlay string, defaultValue int) int {
	data, err := ins.Int(val, overlay)
	if err != nil {
		data = defaultValue
	}
	return data
}

type Config struct {
	data      map[string]interface{}
	separator string
}

func NewConfig(cfgpath string) *Config {
	cfg := new(Config)
	cfg.separator = ":"

	fh, err := os.Open(cfgpath)
	if err != nil {
		cfg.data = make(map[string]interface{})
	} else {
		defer fh.Close()
		confmap, err := ioutil.ReadAll(fh)
		if err != nil {
			panic("Error reading config: " + err.Error())
		}
		if err := yaml.Unmarshal(confmap, &cfg.data); err != nil {
			panic("Error parsing config: " + err.Error())
		}
	}

	return cfg
}

// SetSeparator sets a separator for Find path. Default is ":".
func (cfg *Config) SetSeparator(sep string) *Config {
	cfg.separator = sep
	return cfg
}

// GetContent returns a content of the configuration
func (cfg *Config) Root() *Inspector {
	return NewInspector(&cfg.data)
}

func (cfg *Config) shift(levelPath string, subset map[interface{}]interface{}) map[interface{}]interface{} {
	if subset == nil {
		subset = make(map[interface{}]interface{})
		for k, v := range cfg.data {
			subset[k] = v
		}
	}
	lpath := strings.Split(levelPath, cfg.separator)

	for idx, offset := range lpath {
		idx++
		subset = subset[offset].(map[interface{}]interface{})
		if subset != nil && len(lpath[idx:]) > 0 {
			subset = cfg.shift(strings.Join(lpath[idx:], cfg.separator), subset)
		}
	}
	return subset
}

// Find returns a context of the tree config.
// Each YAML-based config is basically a tree. So Find resets the
// root of the tree to a specific point.
func (cfg *Config) Find(levelPath string) *Inspector {
	subset := make(map[string]interface{})
	for k, v := range cfg.shift(levelPath, nil) {
		subset[k.(string)] = v
	}
	return NewInspector(&subset)
}
