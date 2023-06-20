 # turbo_flag

a drop in replacement for flag package which is included in the core go, but with additional capabilities like 
- Writing command-line apps with subcommands
- Loading configuration file like json,yaml,toml.
- Binding variable/s to values from a configuration file
- Loading `.env` files
- Binding variable/s to environment variable/s
- Enumeration of the values of the flag
- Short alias for a flag

etc.
 
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
// only main command can load the cfg, sub commands / flagset will be preloaded with the same configuration while creating
err = fs.LoadCfg("./test_config/demo.json")
password := fs.String("password", "", "",fs.Cfg("database.password"))

fmt.Println(*password)
//prints "12345"
```

| NOTE: supported config file types:   |
| :------------ |
| *json, yaml/yml, toml.*|

| NOTE:   |
| :------------ |
|  *if you need to ignore flags suplied to your program (i  e act like viper) dont call FlagSet.Parse(args) but rather call FlagSet.ParseWithoutArgs(args)*|



### **binding environment variables**
```go
fs := flag.NewFlagSet("demo", flag.ExitOnError)
//now to bind a flag to a ENV/s
dbPassword:=fs.String("dbPassword","","the password usage",fs.Env( "POSTGRES_PASSWORD", "DB_PASSWORD")) 
//env's set are POSTGRES_PASSWORD=abc
fmt.Println(*dbPassword)
//prints "abc"
``` 
### **setting alias for a flag**
useful for adding a short flag for another flag
```go
fs := flag.NewFlagSet("demo", flag.ExitOnError)
dbPassword:=fs.String("dbPassword","","the password usage string",fs.Alias("p"))
//every property of the original flag will be copied
//when you run the program using "go run . -p "xyz"
fmt.Println(*dbPassword)
//prints "xyz"
```
### **setting enums/options/allowed values for a flag**
```go
//its an error if the default value of a flag is not one of the enums
fs := NewFlagSet("yourProgram", ContinueOnError,fs.Enum("a", "b", "c"))
option := fs.String("option", "c", "")
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

func main() {
	//run our program
	os.Args = []string{"commit", "--branch", "stable"}
	flag.NewMainCmdFs("git", flag.ContinueOnError, os.Args, git)
}

func git(fs *flag.FlagSet, args []string) {
	fs.SubCmdFs("commit", commit)
	fs.SubCmdFs("remote", remote)
	//lets try to commit with branch as argument
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	if branchName != "stable" {
		panic("branchNameFs should be stable")
	}
	//lets run remote with name argument
	err = fs.Parse([]string{"remote", "--name", "origin"})
	if err != nil {
		panic(err)
	}
	if remoteName != "origin" {
		panic("remoteNameFs should be origin")
	}
}

func commit(fs *flag.FlagSet, args []string) {
	var branch string
	fs.StringVar(&branch, "branch", "", "", fs.Alias("b"))

	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	branchName = branch
}

func remote(fs *flag.FlagSet, args []string) {
	var name string
	fs.StringVar(&name, "name", "", "", fs.Alias("n"))

	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	remoteName = name
}

```
## alternative
```go

var (
	branchName = ""
	remoteName = ""
)

func main() {
	//run our program
	os.Args = []string{"commit", "--branch", "stable"}
	flag.NewMainCmd("git", flag.ContinueOnError, os.Args, git)
}

func git(cmd flag.CMD, args []string) {
	cmd.SubCmd("commit", commit)
	cmd.SubCmd("remote", remote)
	//lets try to commit with branch as argument
	err := cmd.Parse(args)
	if err != nil {
		panic(err)
	}
	if branchName != "stable" {
		panic("branchName should be stable")
	}
	//lets run remote with name argument
	err = cmd.Parse([]string{"remote", "--name", "origin"})
	if err != nil {
		panic(err)
	}
	if remoteName != "origin" {
		panic("remoteName should be origin")
	}
}

func commit(fs flag.CMD, args []string) {
	var branch string
	fs.StringVar(&branch, "branch", "", "", fs.Alias("b"))

	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	branchName = branch
}

func remote(fs flag.CMD, args []string) {
	var name string
	fs.StringVar(&name, "name", "", "", fs.Alias("n"))

	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	remoteName = name
}

```