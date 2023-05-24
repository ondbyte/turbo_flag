package flag

/*
import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FlagSet2 struct {
	*flag.FlagSet
	enums      map[string]*Enum
	parseError string
}

// NewFlagSet2 returns a new, empty flag set with the specified name and error handling property. If the name is not empty, it will be printed in the default usage message and in error messages.
func NewFlagSet2(name string, errorHandling ErrorHandling) *FlagSet2 {
	fs := flag.NewFlagSet2(name, flag.ErrorHandling(errorHandling))
	return &FlagSet2{FlagSet2: fs, enums: make(map[string]*Enum)}
}

// parses configuration file to flags available in this flag set
func (fs *FlagSet2) ParseCfgFile(path string) (cfgExists bool, err error) {
	contentMap, notExists, err := fs.tryParseConfig(path)
	if notExists {
		return false, nil
	} else {
		if err == nil {
			err = fs.asignValues(contentMap)
		}
		if err != nil {
			return true, fmt.Errorf("unable to assign values to flags from configuration file\n(%v):\n%v", path, err)
		}
	}
	return true, nil
}

// generates configuration file for this flag set.
// supported configuration file formats include JSON,YAML/YML,PROPERTIES
func (fs *FlagSet2) GenCfgFile(path string) (generated bool, err error) {
	err = fs.tryInitingConfig(path, false)
	if err != nil {
		return false, fmt.Errorf("unable to create configuration file\n(%v):\n%v", path, err)
	}
	return true, nil
}

// asks the user for confirmation before its
// generates configuration file for this flag set.
// supported configuration file formats include JSON,YAML/YML,PROPERTIES
func (fs *FlagSet2) GenCfgFileUserPrompt(path string) (inited bool, err error) {
	err = fs.tryInitingConfig(path, true)
	if err != nil {
		return false, fmt.Errorf("unable to create configuration file\n(%v):\n%v", path, err)
	}
	return true, nil
}

func (fs *FlagSet2) promptForFlags(flags map[string]*flag.Flag) {
	for _, f := range flags {
		enum := fs.enums[f.Name]
		if enum != nil {
			val, err := promptForEnum(enum)
			if err != nil {
				panic(fmt.Sprintf("error getting input from user for enum flag \"%v\": %v", f.Name, err))
			}
			err = f.Value.Set(val)
			if err != nil {
				panic(fmt.Sprintf("error setting val to enum flag \"%v\": %v", f.Name, err))
			}
		} else {
			val, err := promptForValue(f.Usage, func(val string) error {
				return f.Value.Set(val)
			})
			if err != nil {
				panic(fmt.Sprintf("error getting input from user for flag \"%v\": %v", f.Name, err))
			}
			err = f.Value.Set(val)
			if err != nil {
				panic(fmt.Sprintf("error setting val to flag \"%v\": %v", f.Name, err))
			}
		}
	}
}

// prompts for input from user for values of all flags added to this flag set
func (fs *FlagSet2) PromptAll() error {
	flags := make(map[string]*flag.Flag, 0)
	fs.VisitAll(func(f *flag.Flag) {
		flags[f.Name] = f
	})
	fs.promptForFlags(flags)
	return nil
}

// prompts for input from user for values of only flags that were not passed in arguments or parsed ffrom config
func (fs *FlagSet2) PromptForNotPassed() error {
	flags := make(map[string]*flag.Flag, 0)
	fs.VisitAll(func(f *flag.Flag) {
		flags[f.Name] = f
	})
	fs.Visit(func(f *flag.Flag) {
		flags[f.Name] = nil
	})
	fs.promptForFlags(flags)
	return nil
}

func (fs *FlagSet2) tryInitingConfig(cfgFilePath string, askUser bool) error {
	var userResponse string
	if askUser {
		p := newEnum("", fmt.Sprintf("config [%v] does'nt exist,like to init one?", cfgFilePath), "yes", "no")
		r, err := promptForEnum(p)
		if err != nil {
			return fmt.Errorf("unable create config file: " + err.Error())
		}
		userResponse = r
	} else {
		userResponse = "yes"
	}
	if userResponse == "yes" {
		f, err := os.Create(cfgFilePath)
		if err != nil {
			return fmt.Errorf("unable create config file: " + err.Error())
		}
		defer f.Close()
		fileContent := "{"
		i := 0
		fs.VisitAll(func(f *flag.Flag) {
			fileContent += fmt.Sprintf("%v:%v", f.Name, f.DefValue)
			fmt.Println(f.Name)
			i++
		})
		fmt.Println(i)
		fileContent += "}"
		contentMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(fileContent), &contentMap)
		if err != nil {
			return fmt.Errorf("unable create config file: " + err.Error())
		}
		fmt.Println(fileContent)
		fmt.Println(contentMap)
		os.Exit(0)
		contentToWrite := ""
		ext := strings.ToLower(filepath.Ext(cfgFilePath))
		if ext == ".json" {
			contentToWrite, err = WriteMapToJSON(contentMap)
		} else if ext == ".yaml" || ext == ".yml" {
			contentToWrite, err = WriteMapToYAML(contentMap)
		} else if ext == ".properties" {
			contentToWrite = WriteMapToPropertyFile(contentMap)
		} else {
			return fmt.Errorf("[%v] is not supported for config file", ext)
		}
		if err != nil {
			return fmt.Errorf("unable create config file: " + err.Error())
		}
		_, err = f.Write([]byte(contentToWrite))
		if err != nil {
			return fmt.Errorf("unable write to config file: " + err.Error())
		}
	}
	return nil
}

func (fs *FlagSet2) tryParseConfig(cfgFilePath string) (map[string]interface{}, bool, error) {
	b, err := ioutil.ReadFile(cfgFilePath)
	if os.IsNotExist(err) {
		return nil, true, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("unable read config file: " + err.Error())
	}
	contentMap := make(map[string]interface{})
	ext := strings.ToLower(filepath.Ext(cfgFilePath))
	if ext == ".json" {
		contentMap, err = ReadJSONFile(string(b))
	} else if ext == ".yaml" || ext == ".yml" {
		contentMap, err = ReadYAMLFile(string(b))
	} else if ext == ".properties" {
		contentMap, err = ReadPropertyFile(string(b))
	} else {
		return nil, false, fmt.Errorf("[%v] is not supported for config file", ext)
	}
	if err != nil {
		return nil, false, fmt.Errorf("unable read config file: " + err.Error())
	}
	err = fs.asignValues(contentMap)
	if err != nil {
		return nil, false, err
	}
	return contentMap, false, nil
}

func (fs *FlagSet2) asignValues(contentMap map[string]interface{}) error {
	errrS := ""
	fs.VisitAll(func(f *flag.Flag) {
		if val, exists := contentMap[f.Name]; exists {
			v := fmt.Sprint(val)
			err := f.Value.Set(v)
			if err != nil {
				errrS += err.Error()
			}
		}
	})
	if errrS != "" {
		return fmt.Errorf(errrS)
	}
	return nil
}

// EnumString defines a string flag with specified name, default value, and usage string.
// possible values for this flag must passed as options, empty options will result in error
// default value should be any one value from the options otherwise it'll be a error
// The return value is the address of a string variable that stores the value of the flag.
func (fs *FlagSet2) EnumString(name, value, usage string, options ...string) *string {
	f := fs.String(name, value, usage)
	fs.enum(name, value, usage, toDynamicSlice(options))
	return f
}

// Int defines an int flag with specified name, default value, and usage string.
// possible values for this flag must passed as options, empty options will result in error
// default value should be any one value from the options otherwise it'll be a error
// The return value is the address of an int variable that stores the value of the flag.
func (fs *FlagSet2) EnumInt(name string, value int, usage string, options ...int) *int {
	f := fs.Int(name, value, usage)
	fs.enum(name, value, usage, toDynamicSlice(options))
	return f
}

// Float64 defines a float64 flag with specified name, default value, and usage string.
// possible values for this flag must passed as options, empty options will result in error
// default value should be any one value from the options otherwise it'll be a error
// The return value is the address of a float64 variable that stores the value of the flag.
func (fs *FlagSet2) EnumFloat(name string, value float64, usage string, options ...float64) *float64 {
	f := fs.Float64(name, value, usage)
	fs.enum(name, value, usage, toDynamicSlice(options))
	return f
}

func toDynamicSlice[T any](s []T) []interface{} {
	options := make([]interface{}, len(s))
	for i, option := range s {
		options[i] = option
	}
	return options
}

func (fs *FlagSet2) enum(name string, value interface{}, usage string, options []interface{}) {
	defaultIsGood := false
	for _, o := range options {
		if o == value {
			defaultIsGood = true
			break
		}
	}
	if !defaultIsGood {
		e := fmt.Sprintf("enum flag \"%v\" cannot have value \"%v\", only possible values include [%v]\n", name, value, stringS(options))
		e += fmt.Sprintf("use normal flag rather than enum flag if you want to get any possible values as input, or kindly add \"%v\" to enum options\n", value)
		panic(color(e, Red))
	}
	p := &Enum{
		name:    name,
		usage:   usage,
		options: options,
	}
	fs.enums[name] = p
}
*/
