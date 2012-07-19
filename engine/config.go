package engine

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	FileName string
	values   map[string]interface{}
}

func NewCfg(fileName string) (*Config, error) {
	cfg := new(Config)

	cfg.FileName = fileName
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
	cfg, err := NewCfg("excavation.cfg")
	if err != nil {
		return nil, err
	}

	if file, err := os.Open(cfg.FileName); err != nil {
		if os.IsNotExist(err) {
			//file doesn't exist
			// create one with default values
			defaultConfigHandler(cfg)
			/*
				cfg.SetValue("WindowWidth", 1024)
				cfg.SetValue("WindowHeight", 728)
				cfg.SetValue("WindowDepth", 24)
				cfg.SetValue("Fullscreen", false)
				cfg.SetValue("VSync", 1)
			*/
			if err = cfg.Write(); err != nil {
				return nil, err
			}
		}
		return nil, err
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
			//	Mouse_AxisX
			//	Joy4_Axis3
			defaultConfigHandler(cfg)
			if err = cfg.Write(); err != nil {
				return nil, err
			}
		}
		return nil, err
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
	return cfg.values[name]
}

func (cfg *Config) Int(name string) int {
	return int(cfg.values[name].(float64))
}

func (cfg *Config) String(name string) string {
	return cfg.values[name].(string)
}

func (cfg *Config) Bool(name string) bool {
	return cfg.values[name].(bool)
}

func (cfg *Config) Float(name string) float64 {
	return cfg.values[name].(float64)
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
