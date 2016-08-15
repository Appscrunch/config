package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/kardianos/osext"
)

var configPath string

//for logger
var CallerInfo = "false"

// ReadGlobalConfig unmarshal json-object cf
// If parsing was not successuful, function return a structure with default options
func ReadGlobalConfig(cf interface{}, whatParsed string) error {

	file, e := ioutil.ReadFile(GetConfigFilename())
	if e != nil {
		//log.WithCaller(slf.CallerShort).Errorf("Error: %s\n", e.Error())
		fmt.Fprintf(os.Stderr, "[config] Error: %s\n", e.Error())
		return e
	}

	if err := json.Unmarshal([]byte(file), cf); err != nil {
		//log.WithCaller(slf.CallerShort).Errorf("Error parsing JSON: %s. For [%s] will be used defaulf options.",
		//	err.Error(), whatParsed)
		fmt.Fprintf(os.Stderr, "[config] Error: %s\n", e.Error())
		return err
	}
	fmt.Fprintf(os.Stderr, "[config] Parsed [%s] configuration from [%s] file.\n", whatParsed, GetConfigFilename())
	//fmt.Fprintf(os.Stderr, "[config] If a field has wrong format, will be used default value.\n")

	//fmt.Printf("%v\n", cf)
	return nil
}

// GetConfigFilename is a function fot getting a name of a binary with full path to it
func GetConfigFilename() string {

	if configPath != "" {
		return configPath
	}

	binaryPath, err := osext.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not get a path to binary file: %s\n", err.Error())
	}
	if runtime.GOOS == "windows" {
		// without ".exe"
		//TODO: FIX use path func
		binaryPath = binaryPath[:len(binaryPath)-4]
		fmt.Fprintf(os.Stderr, "[config] Config file for windows %s\n", binaryPath)
	}

	configPath = binaryPath + ".config"
	return configPath
}
