package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	//"flag"
	"reflect"
	"runtime"

	"github.com/kardianos/osext"
)

var configPath string

//for logger
var CallerInfo = "false"

type configConfigT struct {
	Verbose  bool
	VVerbose bool
}

var cfg = configConfigT{}

func init() {
	//read command line before config loading
	readOsArgs(&cfg)
	//read config
	ReadGlobalConfig(&cfg, "default config")
	/*for _,tmp := range os.Args {
		fmt.Println(tmp)
	}*/
}

type fieldDesc struct {
	name  string
	value reflect.Value
	field reflect.StructField
}

//read all fields and values from structure 'v' into slice 'fields'
func flatFieldsRecursive(v reflect.Value, prefix string, fields *map[string]fieldDesc) {
	if cfg.VVerbose {
		fmt.Fprintf(os.Stderr, "flatFieldsRecursive %v len:%d\n", v, len(*fields))
	}

	var tmpPrefix = ""
	if prefix != "" {
		tmpPrefix = prefix + "."
	}
	for i := 0; i < v.NumField(); i++ {
		valueField := v.Field(i)
		fld := v.Type().Field(i)
		if fld.Type.Kind() == reflect.Struct {
			flatFieldsRecursive(valueField, tmpPrefix+fld.Name, fields)
		} else {
			(*fields)[tmpPrefix+fld.Name] = fieldDesc{
				name:  tmpPrefix + fld.Name,
				field: fld,
				value: valueField,
			}
		}
	}
}

//create recursive list of fields and values in struct v
func flatFields(v reflect.Value) *map[string]fieldDesc {
	//fmt.Println("flatFields", t)
	result := make(map[string]fieldDesc)
	pResult := &result
	flatFieldsRecursive(v, "", pResult)
	return pResult
}

func setAsInt(v reflect.Value, s string, bitSize int) error {
	i, err := strconv.ParseInt(s, 10, bitSize)
	if err != nil {
		return fmt.Errorf("Error parse field %s value %s as int%d: %s", v, s, bitSize, err)
	}
	v.SetInt(i)
	return nil
}

func setAsFloat(v reflect.Value, s string, bitSize int) error {
	f, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		return fmt.Errorf("Error parse field %s value %s as float%d: %s", v, s, bitSize, err)
	}
	v.SetFloat(f)
	return nil
}

func setAnyAsString(v reflect.Value, s string) error {
	if !v.CanSet() {
		return fmt.Errorf("Cannot set value %v", v)
	}
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
		return nil
	case reflect.Int:
		return setAsInt(v, s, 0)
	case reflect.Int8:
		return setAsInt(v, s, 8)
	case reflect.Int16:
		return setAsInt(v, s, 16)
	case reflect.Int32:
		return setAsInt(v, s, 32)
	case reflect.Int64:
		return setAsInt(v, s, 64)
	case reflect.Bool:
		ss := strings.ToLower(s)
		if strings.HasPrefix(ss, "y") {
			v.SetBool(true)
			return nil
		}
		if strings.HasPrefix(ss, "n") {
			v.SetBool(true)
			return nil
		}
		b, err := strconv.ParseBool(ss)
		if err != nil {
			return fmt.Errorf("Error parse field %s value %s as bool: %s", v, s, err)
		}
		v.SetBool(b)
		return nil
	case reflect.Float32:
		return setAsFloat(v, s, 32)
	case reflect.Float64:
		return setAsFloat(v, s, 64)
	}
	return fmt.Errorf("Error parse field %s value %s. Unknown kind %v ", v, s, v.Kind())
}

func readOsArgsInner(cf interface{}) error {
	rv := reflect.ValueOf(cf)

	if rv.Kind() != reflect.Ptr || rv.IsNil() || reflect.Indirect(rv).Kind() != reflect.Struct {
		return fmt.Errorf("Invalid config type [%v] [%v] [%v], should be pointer to struct", reflect.TypeOf(rv), rv.Kind(), reflect.Indirect(rv).Kind())
	}
	value := reflect.Indirect(rv)
	fields := flatFields(value)

	/*if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "[config] dump fields\n")
		for _,fld := range *fields {
			fmt.Fprintf(os.Stderr, "[config] %s : %v = \"%v\"\n", fld.name, fld.field.Type, fld.value)
		}
	}*/

	for _, arg := range os.Args[1:] {
		//fmt.Println(arg)
		if arg == "--" {
			break
		}
		if !strings.HasPrefix(arg, "--") {
			continue
		}
		kv := strings.SplitN(arg, "=", 2)
		//fmt.Println(kv)
		if len(kv) < 1 {
			continue
		}
		key := strings.TrimSpace(kv[0])[2:]

		fld, exists := (*fields)[key]
		if !exists {
			if cfg.VVerbose {
				fmt.Fprintf(os.Stderr, "Field not found %v\n", key)
			}
			continue
		}

		if len(kv) >= 2 {
			err := setAnyAsString(fld.value, kv[1])
			if err != nil && cfg.Verbose {
				fmt.Fprintf(os.Stderr, "error parse field %s\n", err)
			}
		} else {
			fld.value.Set(reflect.ValueOf(true))
		}
	}

	return nil
}

func dumpFields(cf interface{}, whatParsed string) {
	rv := reflect.ValueOf(cf)

	if rv.Kind() != reflect.Ptr || rv.IsNil() || reflect.Indirect(rv).Kind() != reflect.Struct {
		return
	}
	value := reflect.Indirect(rv)
	//fmt.Println("ReadGlobalConfig", rv, rv.Kind(), reflect.Indirect(rv).Kind())
	fields := flatFields(value)

	//if cfg.verbose {
	fmt.Fprintf(os.Stderr, "[config] dump fields for %s\n", whatParsed)
	for _, fld := range *fields {
		fmt.Fprintf(os.Stderr, "[config] %s : %v = \"%s\"\n", fld.name, fld.field.Type, fld.value)
	}
	//}
}

func readOsArgs(cf interface{}) {
	err := readOsArgsInner(cf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] [config] parse os flags [%v]", err)
	}
}

// ReadGlobalConfig unmarshal json-object cf
// If parsing was not successuful, function return a structure with default options
func readGlobalConfigInner(cf interface{}, filename string) error {

	file, e := ioutil.ReadFile(filename)
	if e != nil {
		//log.WithCaller(slf.CallerShort).Errorf("Error: %s\n", e.Error())
		//fmt.Fprintf(os.Stderr, "[config] Error: %s\n", e.Error())
		return e
	}

	if err := json.Unmarshal([]byte(file), cf); err != nil {
		//log.WithCaller(slf.CallerShort).Errorf("Error parsing JSON: %s. For [%s] will be used defaulf options.",
		//	err.Error(), whatParsed)
		//fmt.Fprintf(os.Stderr, "[config] Error: %s\n", e.Error())
		return err
	}
	//fmt.Fprintf(os.Stderr, "[config] Parsed [%s] configuration from [%s] file.\n", whatParsed, GetConfigFilename())
	//fmt.Fprintf(os.Stderr, "[config] If a field has wrong format, will be used default value.\n")

	//fmt.Printf("%v\n", cf)
	return nil
}

func ReadGlobalConfig(cf interface{}, whatParsed string) {
	filename := GetConfigFilename()
	err := readGlobalConfigInner(cf, filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] [config] parse [%s] error [%v]", whatParsed, err)
	} else {
		if cfg.Verbose {
			fmt.Fprintf(os.Stderr, "[INFO]  [config] Parsed [%s] configuration from file [%s].\n", whatParsed, filename)
		}
	}
	readOsArgs(cf)
	if cfg.Verbose {
		dumpFields(cf, whatParsed)
	}
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
		if cfg.Verbose {
			fmt.Fprintf(os.Stderr, "[config] Config file for windows %s\n", binaryPath)
		}
	}

	configPath = binaryPath + ".config"
	return configPath
}
