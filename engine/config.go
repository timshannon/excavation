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
	Name     string
	FileName string
	values   map[string]interface{}
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
			//TODO: Default controls
			// Naming will be as follows
			//  Device_Axis#
			//  Device_Button#
			//  Examples:
			//	Joy12_13
			//	Key_W
			//	Mouse_Axis1
			//	Joy4_Axis3
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
		return errors.New("No Filename set for Config object")
	}

	data, err := ioutil.ReadFile(cfg.FileName)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &cfg.values); err != nil {
		return err
	}
	return nil
}

func (cfg *Config) Value(name string) interface{} {
	if cfg.values[name] == nil {
		return nil
	}
	return cfg.values[name]
}

func (cfg *Config) Int(name string) int {
	return int(cfg.Value(name).(float64))
}

func (cfg *Config) String(name string) string {
	return cfg.Value(name).(string)
}

func (cfg *Config) Bool(name string) bool {
	return cfg.Value(name).(bool)
}

func (cfg *Config) Float(name string) float64 {
	return cfg.Value(name).(float64)
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
