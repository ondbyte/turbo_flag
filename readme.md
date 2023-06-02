 # turbo_flag

a drop in replacement for flag package which is included in the core go, but with additional capabilities like loading values to flags from a configuration file, environment variable, enumeration of the values of the flag, sub commands, short alias for a flag etc.

**A wannabe viper or cobra alternative.**
 
| NOTE: following functions/methods has been deprecated because they were not adding any value to the working of the core flag package:   |
| :------------ |
| ~~FlagSet.PrintDefaults()~~: was an alias for usage itself|
| ~~FlagSet.Usage()~~: rather you can use GetDefaultUsage() which will return a string to print to the console|
|~~ErrHelp~~: not needed|
|~~FlagSet.outPut~~|
|~~FlagSet.SetOutput(io.Writer)~~: not needed|
|automatic handling of \--help or -h flag has been removed|

### drop in replacement to core flag

```go
import "flag"

fs:=flag.NewFlagSet("demo",flag.ContinueOnError)
help:=fs.Bool("help",false,"prints help")
```
just change the import to "github.com/ondbyte/turbo_flag"
```go
import "github.com/ondbyte/turbo_flag"

fs:=flag.NewFlagSet("demo",flag.ContinueOnError)
help:=fs.Bool("help",false,"prints help")
```

## what additional features it has over flag

### **loading configurations**

lets load a file called demo.json having content
```json
{
    "database":{
        "password":"12345"
    }
}
```
now bind the cfg
```go
fs := NewFlagSet("test", ContinueOnError)
password := fs.String("password", "", "")
err = fs.LoadCfg("./test_config/demo.json")
//pass the password ptr to BindCfg, use dot notation to access the value in the cfg
err = fs.BindCfg(password, "database.password")
fmt.Println(*password)
//prints "12345"
```

| NOTE: supported config file types:   |
| :------------ |
| *json, yaml/yml, toml.*|

| NOTE:   |
| :------------ |
|  *if you need to ignore flags suplied your program (i  e act like viper) dont call FlagSet.Parse(args) but rather call FlagSet.ParseWithoutArgs(args)*|

### **generating configurtion file** (_your own git init_)


### **binding environment variables**
```go
fs := flag.NewFlagSet("demo", flag.ExitOnError)
//now to bind a flag to a ENV/s
dbPassword:=fs.String("dbPassword","","the password usage") 
err = fs.BindEnv(dbPassword, "POSTGRES_PASSWORD", "DB_PASSWORD")
if err!=nil{
	//failed to bind
}
//env's set are POSTGRES_PASSWORD=abc
fmt.Println(*dbPassword)
//prints "abc"
``` 
### **setting alias for flags**


```go
fs := flag.NewFlagSet("demo", flag.ExitOnError)
dbPassword:=fs.String("dbPassword","","the password usage string")
err := fs.Alias(dbPassword, "p")
//every property of the original flag will be copied
//when you run the program using "go run . -p "xyz"
fmt.Println(*dbPassword)
//prints "xyz"
```
### **setting enums/options/allowed values for a flag**
```go
fs := NewFlagSet("yourProgram", ContinueOnError)
option := fs.String("option", "c", "")
//its an error if the default value of a flag is not one of the enums
err := fs.BindEnum(option, "a", "b", "c")
//from the commandline  allowed values are
//yourProgram -option a
//yourProgram -option b
//yourProgram -option c
// otherwise its a error
```
### **Sub-commands**
example: _a git program with commit and remote sub-commands_
```go
var (
	branchName = ""
	remoteName = ""
)
func git() {
	fs := flag.NewFlagSet("git", flag.ContinueOnError)
	fs.SubCmd("commit", commit)
	fs.SubCmd("remote", remote)
	//lets try to commit with branch as argument
	err := fs.Parse([]string{"commit", "--branch", "stable"})
	if err != nil {
		panic(err)
	}
	if branchName != "stable" {
		panic("branchName should be stable")
	}
	//lets run remote with name argument
	err = fs.Parse([]string{"remote", "--name", "origin"})
	if err != nil {
		panic(err)
	}
	if remoteName != "origin" {
		panic("remoteName should be origin")
	}
}

func commit(args []string) {
	fs := flag.NewFlagSet("commit", flag.ContinueOnError)
	var branch string
	fs.StringVar(&branch, "branch", "", "")
	fs.Alias(&branch, "b")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	branchName = branch
}

func remote(args []string) {
	fs := flag.NewFlagSet("remote", flag.ContinueOnError)
	var name string
	fs.StringVar(&name, "name", "", "")
	fs.Alias(&name, "n")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	remoteName = name
}
```