//MIT License
//Copyright (c) [2017] [dzthink]

//配置解析和访问
package dogo

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"strings"
	"errors"
)

type Config struct {
	conf map[string]*json.RawMessage
}

func NewConfig(path string) (*Config, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf map[string]*json.RawMessage
	err = json.Unmarshal(dat, &conf)
	if err != nil {
		return nil, err
	}

	c := &Config{
		conf : conf,
	}
	c.imports()
	return c, nil
}

func(c *Config)imports() {
	if importRaw, ok := c.conf["@imports"]; ok {
		var imports map[string]string
		err := json.Unmarshal(*importRaw, &imports)
		if err == nil {
			for k, importPath := range imports {
				dat, err := ioutil.ReadFile(importPath)
				if err != nil {
					continue
				}
				childConfRaw := json.RawMessage(dat)
				c.conf[k] = &childConfRaw
			}
		}
		delete(c.conf, "@imports")
	}

}

func(c *Config)Child(key string) (*Config, error) {
	if raw, err := c.fetchKey(key); err == nil {
		var conf map[string]*json.RawMessage
		err := json.Unmarshal(*raw, &conf)
		if err != nil {
			return nil, err
		}
		return &Config{
			conf : conf,
		}, nil
	}
	return nil, fmt.Errorf("config item %s not exist", key)
}

func(c *Config)ChildList(key string)([]*Config, error) {
	if raw, err := c.fetchKey(key); err == nil {
		var confs []map[string]*json.RawMessage
		err := json.Unmarshal(*raw, &confs)
		if err != nil {
			return nil, err
		}
		configs := make([]*Config, 0, len(confs))
		for _, v := range confs {
			configs = append(configs, &Config{
				conf : v,
			})
		}
		return configs, nil
	}
	return nil,fmt.Errorf("config item %s not exist", key)
}

func(c *Config)Bool(key string)(bool, error) {
	if raw, err := c.fetchKey(key); err == nil {
		var b bool
		if err := json.Unmarshal(*raw, &b); err != nil {
			return false, err
		}
		return b, nil
	}
	return false, fmt.Errorf("config item %s not exist", key)
}

func(c *Config)Int(key string) (int64 ,error) {
	if raw, err := c.fetchKey(key); err == nil {
		var b int64
		if err := json.Unmarshal(*raw, &b); err != nil {
			return 0, err
		}
		return b, nil
	}
	return 0, fmt.Errorf("config item %s not exist", key)
}

func(c *Config)String(key string) (string, error) {
	if raw, err := c.fetchKey(key); err == nil {
		var b string
		if err := json.Unmarshal(*raw, &b); err != nil {
			return "", err
		}
		return b, nil
	}
	return "", fmt.Errorf("config item %s not exist", key)
}

func(c *Config)Float(key string) (float64, error) {
	if raw, err := c.fetchKey(key); err == nil {
		var b float64
		if err := json.Unmarshal(*raw, &b); err != nil {
			return 0.0, err
		}
		return b, nil
	}
	return 0.0, fmt.Errorf("config item %s not exist", key)
}

func(c *Config)Get(key string, v interface{}) (error) {
	if raw, err := c.fetchKey(key); err == nil {
		if err := json.Unmarshal(*raw, v); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("config item %s not exist", key)
}

func(c *Config)fetchKey(key string) (*json.RawMessage, error) {
	if raw, ok := c.conf[key]; ok {
		return raw, nil
	}
	fmt.Println(key)
	keys := strings.Split(key, ".")
	if len(keys) <= 1 {
		return nil, errors.New("config " + key + "not found")
	}
	if cConf, err := c.Child(keys[0]); err != nil {
		return nil , err
	} else {
		return cConf.fetchKey(strings.Join(keys[1:], "."))
	}
}

func(c *Config)ToString() string {
	dat, err := json.Marshal(c.conf)
	if err != nil {
		return ""
	}
	return string(dat)
}

