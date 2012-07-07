package engine

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type Config map[string]interface{}

//Loads the standard config file 
func LoadConfig() (Config, error) {
	var cfgFile string
	userDir, err := UserDir()
	if err != nil {
		return nil, err
	}

	cfgFile = path.Join(userDir, "config.cfg")
	if file, err := os.Open(cfgFile); err != nil {
		if os.IsNotExist(err) {
			//file doesn't exist
			// create one with default values
			createDefaultConfig(cfgFile)
		}
	} else {
		file.Close()
	}
	return LoadConfigFile(cfgFile)
}

func LoadControlConfig() (Config, error) {
	userDir, err := UserDir()
	if err != nil {
		return nil, err
	}
	//TODO: Control mapping
	return LoadConfigFile(path.Join(userDir, "controls.cfg"))

}

//Loads a specific config file at a specific location
func LoadConfigFile(file string) (Config, error) {
	cfg := Config{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg Config) Value(name string) interface{} {
	return cfg[name]
}

func (cfg Config) Int(name string) int {
	return int(cfg[name].(float64))
}

func (cfg Config) String(name string) string {
	return cfg[name].(string)
}

func (cfg Config) Bool(name string) bool {
	return cfg[name].(bool)
}

func (cfg Config) Float(name string) float64 {
	return cfg[name].(float64)
}

func createDefaultConfig(file string) error {
	cfg := Config{}
	cfg["WindowWidth"] = 1024
	cfg["WindowHeight"] = 728
	cfg["WindowDepth"] = 24
	cfg["Fullscreen"] = false
	cfg["VSync"] = 1

	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, data, 0644)
	return err
}
