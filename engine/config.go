// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

package engine

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	Name           string
	FileName       string
	values         map[string]interface{}
	onWriteHandler ConfigOnWriteHandler
}

func NewCfg(fileName string) (*Config, error) {
	cfg := new(Config)
	cfg.Name = fileName
	cfg.values = make(map[string]interface{})
	//if just a filename with no path is passed in,
	// then combine it with the userDir
	if !path.IsAbs(fileName) {
		userDir, err := UserDir()
		if err != nil {
			return nil, err
		}
		fileName = path.Join(userDir, fileName)
	}

	cfg.FileName = fileName
	return cfg, nil
}

//Loads the standard config file 
func NewStandardCfg() (*Config, error) {
	cfg, err := NewCfg(appName + ".cfg")
	if err != nil {
		return nil, err
	}

	if file, err := os.Open(cfg.FileName); err != nil {
		if os.IsNotExist(err) {
			//file doesn't exist
			// create one with default values
			defaultConfigHandler(cfg)
			if err = cfg.Write(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		file.Close()
	}
	return cfg, nil
}

func NewControlCfg() (*Config, error) {
	cfg, err := NewCfg("controls.cfg")
	if err != nil {
		return nil, err
	}

	if file, err := os.Open(cfg.FileName); err != nil {
		if os.IsNotExist(err) {
			//file doesn't exist
			// create one with default values
			defaultConfigHandler(cfg)
			if err = cfg.Write(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		file.Close()
	}
	return cfg, nil

}

type DefaultConfigHandler func(cfg *Config)

var defaultConfigHandler DefaultConfigHandler

//SetDefaultConfigHandler lets the game code define the
// defaults that should be written if a config file doesn't exist
func SetDefaultConfigHandler(function DefaultConfigHandler) {
	defaultConfigHandler = function
}

//Loads a specific config file at a specific location
func (cfg *Config) Load() error {
	if cfg.FileName == "" {
		err := errors.New("No Filename set for Config object")
		RaiseError(err)
		return err
	}

	data, err := ioutil.ReadFile(cfg.FileName)
	if err != nil {
		RaiseError(err)
		return err
	}
	if err = json.Unmarshal(data, &cfg.values); err != nil {
		RaiseError(err)
		return err
	}
	return nil
}

func (cfg *Config) Value(name string) interface{} {
	return cfg.values[name]
}

//DefaultValueHandler is a function  you can define to
// return the a default value if a given configuration entry isn't found
// in a config file.  It allows you to set sane defaults
type DefaultValueHandler func(name string) interface{}

var defaultValueHandlers = make([]DefaultValueHandler, 10)

//SetDefaultConfigHandler lets the game code define the
// defaults that should be written if a config file doesn't exist
func SetDefaultValueHandler(function DefaultValueHandler) {
	defaultValueHandlers = append(defaultValueHandlers, function)
}

func handleMissing(name string) (value interface{}) {
	RaiseError(errors.New("Config entry " + name + " does not exist. Using default."))
	for i := range defaultValueHandlers {
		value = defaultValueHandlers[i]
		if value != nil {
			return value
		}
	}

	return nil
}

func (cfg *Config) Int(name string) int {
	value, ok := cfg.values[name].(float64)
	if !ok {
		value := handleMissing(name)
		if value != nil {
			return value.(int)
		}

		return 0
	}
	return int(value)
}

func (cfg *Config) String(name string) string {
	value, ok := cfg.values[name].(string)
	if !ok {
		value := handleMissing(name)
		if value != nil {
			return value.(string)
		}
		return ""
	}
	return value
}

func (cfg *Config) Bool(name string) bool {
	value, ok := cfg.values[name].(bool)
	if !ok {
		value := handleMissing(name)
		if value != nil {
			return value.(bool)
		}
		return false
	}
	return value
}

func (cfg *Config) Float(name string) float32 {
	value, ok := cfg.values[name].(float64)
	if !ok {
		value := handleMissing(name)
		if value != nil {
			return value.(float32)
		}
		return 0.0
	}
	return float32(value)
}

func (cfg *Config) SetValue(name string, value interface{}) {
	cfg.values[name] = value
}

func (cfg *Config) Write() error {
	data, err := json.MarshalIndent(cfg.values, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(cfg.FileName, data, 0644)

	return err

}

type ConfigOnWriteHandler func(cfg *Config)

//RegisterOnWriteHandler registers a function to be called when
// this config file is written.  So that if changes are made, 
// the consumers of the config can get the latest values
func (cfg *Config) RegisterOnWriteHandler(handler ConfigOnWriteHandler) {
	cfg.onWriteHandler = handler
}
