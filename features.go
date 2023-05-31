package flag

import (
	"fmt"
	"os"
	"sort"
)

type orderOfFlagFeature int
type orderOfFlagSetFeature int

const (
	_ orderOfFlagFeature = iota
	orderOfBindEnv
	orderOfBindConfig
	orderOfBindFlag
)
const (
	_ orderOfFlagSetFeature = iota
	orderOfConfig
	orderOfSubCmd
)

type flagFeature struct {
	order   orderOfFlagFeature
	flag    *Flag
	flagSet *FlagSet
	add     func(*FlagSet, *Flag) error
}

type flagFeatures []*flagFeature

func (a flagFeatures) Len() int           { return len(a) }
func (a flagFeatures) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a flagFeatures) Less(i, j int) bool { return a[i].order < a[j].order }

//we bind flags last, bind config before that, bind env before doing that and so on
//order will be like this cmd args > configuration > env etc
//so we sort flag features
func sortFlagFeatures(features []*flagFeature) {
	ffs := flagFeatures(features)
	sort.Stable(ffs)
}

func isEnumValid(e string, enums []string) bool {
	valid := false
	if len(enums) > 0 {
		for _, s := range enums {
			if e == s {
				valid = true
			}
		}
	} else {
		valid = true
	}
	return valid
}

//possible values for the flag, if you add this feature, flag can only have the values defined in ths enums
func EnumsOld(enums ...string) *flagFeature {
	return &flagFeature{
		order: 0,
		add: func(fs *FlagSet, ff *Flag) error {
			for _, enum := range enums {
				ff.enums[enum] = true
			}
			if !isEnumValid(ff.DefValue, enums) {
				return fmt.Errorf("default value %v doesnt exists in the enums %v for flag name %v", ff.DefValue, enums, ff.Name)
			}
			return nil
		},
	}
}

//alias for BindFlag
func Short(name string) *flagFeature {
	return BindFlagOld(name)
}

//binds another flag with name for the flag you're defining,
//you dont need to add a new flag to the flagset,new flag will be created for you,
//useful for defining short name for the flag you're defining or multiple flag names and single value ptr,
//bindings will also be copied
func BindFlagOld(names ...string) *flagFeature {
	return &flagFeature{
		order: orderOfBindFlag,
		add: func(fs *FlagSet, ff *Flag) error {
			for _, name := range names {
				if name == ff.Name {
					return fmt.Errorf("cannot bind to the flag with the same name %v", name)
				}
				_, err := fs.varErr(ff.Value, name, ff.Usage)
				if err != nil {
					return fmt.Errorf("error while binding flag %v to flag %v : %v", name, ff.Name, err)
				}
			}
			return nil
		},
	}
}

//binds environment variables to the flag you are defining,
//last non empty value will be assigned,
//passing argument value on the commandline will override the value of this flag even if the env is set
func BindEnvOld(envs ...string) *flagFeature {
	return &flagFeature{
		order: orderOfBindEnv,
		add: func(fs *FlagSet, ff *Flag) error {
			for _, env := range envs {
				ff.envs[env] = true
				val := os.Getenv(env)
				if val != "" {
					err := ff.Set(val)
					if err != nil {
						return fmt.Errorf("error while setting value from environment, flag name %v,env %v,value %v : %v", ff.Name, env, val, err)
					}
				}
			}
			return nil
		},
	}
}

//binds this flag to value from config file
//use dot notation
//multiple values are supported last found value will be loaded
func BindCfgOlds(cfgs ...string) *flagFeature {
	return &flagFeature{
		order: orderOfBindConfig,
		add: func(fs *FlagSet, ff *Flag) error {
			if fs.cfg == nil {
				return fmt.Errorf(fmt.Sprint(`you are binding the flag `,
					ff.Name,
					` to config notation/s `,
					cfgs,
					` but you've not loaded any config file, load the config using flag.Config feature while calling the flag.NewFlagSet`,
					` example:`,
					"\n\n",
					` fs := NewFlagSet("example", ExitOnError, Config("./awesome.json"))`,
					"\n\n"))
			}
			for _, notation := range cfgs {
				val, err := getValueByDotNotation(fs.cfg, notation)
				if err == nil {
					err := ff.Set(val)
					if err != nil {
						return fmt.Errorf("unable to set notation %v value %v to flag %v", notation, val, ff.Name)
					}
				}
			}
			for _, cfg := range cfgs {
				ff.cfgs[cfg] = true
			}
			return nil
		},
	}
}
