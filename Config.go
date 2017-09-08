//MIT License
//Copyright (c) [2017] [dzthink]

//配置解析和访问
package dogo

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
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
	err = json.Unmarshal(dat, conf)
	if err != nil {
		return nil, err
	}
	return &Config{
		conf : conf,
	}, nil
}

func(c *Config)Child(key string) (*Config, error) {
	if raw, ok := c.conf[key]; ok {
		var conf map[string]*json.RawMessage
		err := json.Unmarshal(*raw, conf)
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
	if raw, ok := c.conf[key]; ok {
		var confs []map[string]*json.RawMessage
		err := json.Unmarshal(*raw, confs)
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
	if raw, ok := c.conf[key]; ok {
		var b bool
		if err := json.Unmarshal(*raw, &b); err != nil {
			return false, err
		}
		return b, nil
	}
	return false, fmt.Errorf("config item %s not exist", key)
}

func(c *Config)Int(key string) (int64 ,error) {
	if raw, ok := c.conf[key]; ok {
		var b int64
		if err := json.Unmarshal(*raw, &b); err != nil {
			return 0, err
		}
		return b, nil
	}
	return 0, fmt.Errorf("config item %s not exist", key)
}

func(c *Config)String(key string) (string, error) {
	if raw, ok := c.conf[key]; ok {
		var b string
		if err := json.Unmarshal(*raw, &b); err != nil {
			return "", err
		}
		return b, nil
	}
	return "", fmt.Errorf("config item %s not exist", key)
}

func(c *Config)Float(key string) (float64, error) {
	if raw, ok := c.conf[key]; ok {
		var b float64
		if err := json.Unmarshal(*raw, &b); err != nil {
			return 0.0, err
		}
		return b, nil
	}
	return 0.0, fmt.Errorf("config item %s not exist", key)
}

func(c *Config)Get(key string, v interface{}) (interface{}, error) {
	if raw, ok := c.conf[key]; ok {
		if err := json.Unmarshal(*raw, v); err != nil {
			return 0.0, err
		}
		return v, nil
	}
	return nil, fmt.Errorf("config item %s not exist", key)
}

