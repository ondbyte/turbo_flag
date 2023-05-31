# turbo_flag

a drop in replacement for flag package that comes with go which is a package to write commandline apps, with additional capabilties like loading values to flags from a configuration file, environement variable, enumeration or even binding to another flag.

### drop in replacement

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

### what additional features it has over flag

**loading configurations**

| note:  |
| :------------ |
| supported config file types:  json, yaml/yml, properties|
**lets load a file called demo.json having content**
```json
{
    "database":{
        "password":"12345"
    }
}
```
now bind the cfg
```go
fs := flag.NewFlagSet("demo", flag.ExitOnError, flag.Config("./demo.json"))
//now to bind a flag to a value in the config file loaded
dbPassword:=fs.String("dbPassword","","the password usage string",flag.BindCfg("database.password")) //use dot notation to access the value in the cfg
fmt.Println(*dbPassword)
//prints "12345"
```
**binding environment variables**
```go
fs := flag.NewFlagSet("demo", flag.ExitOnError)
//now to bind a flag to a ENV/s
dbPassword:=fs.String("dbPassword","","the password usage string",flag.BindEnv("DB_PASSWORD","POSTGRES_PASSWORD")) 

//env set POSTGRES_PASSWORD=abc
fmt.Println(*dbPassword)
//prints "abc"
```
**setting short flags**

```go
fs := flag.NewFlagSet("demo", flag.ExitOnError)
dbPassword:=fs.String("dbPassword","","the password usage string",flag.Short("p")) 
//when you run the program using "go run . -p "xyz"
fmt.Println(*dbPassword)
//prints "xyz"
```