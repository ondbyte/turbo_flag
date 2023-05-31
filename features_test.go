package flag

import (
	"fmt"
	"os"
	"testing"
)

func TestWithoutConfig(t *testing.T) {
	fs := NewFlagSet("example", ExitOnError)

	fs.String("password", "", "password", BindFlag("p"), BindCfg("database.password"))
	err := fs.Parse(make([]string, 0))
	if err == nil {
		t.Errorf("expected error")
	}
	fmt.Println(err)
}
func TestWithConfigWhichDoesntExist(t *testing.T) {
	fs := NewFlagSet("example", ExitOnError)
	err := fs.LoadCfg("./cfgxyzsdfc.json")
	if err == nil {
		t.Errorf("expected error")
	}
	fs.String("password", "", "password", BindFlag("p"), BindCfg("database.password"))
	err = fs.Parse(make([]string, 0))
	if err == nil {
		t.Errorf("expected error")
	}
	fmt.Println(err)
}
func TestWithConfigWithUnsupportedExtension(t *testing.T) {
	fs := NewFlagSet("example", ExitOnError)

	err := fs.LoadCfg("./testing/demo.abc")
	if err == nil {
		t.Errorf("expected error")
	}
	fs.String("password", "", "password", BindFlag("p"), BindCfg("database.password"))
	err = fs.Parse(make([]string, 0))
	if err == nil {
		t.Errorf("expected error")
	}
	fmt.Println(err)
}
func TestWithConfigJSON(t *testing.T) {
	fs := NewFlagSet("example", ExitOnError)
	err := fs.LoadCfg("./testing/demo.json")
	if err != nil {
		t.Fatal(err)
	}
	p := fs.String("password", "", "password", BindFlag("p"), BindCfg("database.password"))
	err = fs.Parse(make([]string, 0))
	if err != nil {
		t.Fatal(err)
	}
	if *p != "12345" {
		t.Errorf("expected password to be 12345 but %v", *p)
	}
}
func TestWithConfigYAML(t *testing.T) {
	fs := NewFlagSet("example", ExitOnError)
	err := fs.LoadCfg("./testing/demo.yaml")
	if err != nil {
		t.Fatal(err)
	}
	p := fs.String("password", "", "password", BindFlag("p"), BindCfg("database.password"))
	err = fs.Parse(make([]string, 0))
	if err != nil {
		t.Fatal(err)
	}
	if *p != "12345" {
		t.Errorf("expected password to be 12345 but %v", *p)
	}
}
func TestWithConfigTOML(t *testing.T) {
	fs := NewFlagSet("example", ExitOnError)
	err := fs.LoadCfg("./testing/demo.toml")
	if err != nil {
		t.Fatal(err)
	}
	p := fs.String("password", "", "password", BindFlag("p"), BindCfg("database.password"))
	err = fs.Parse(make([]string, 0))
	if err != nil {
		t.Fatal(err)
	}
	if *p != "12345" {
		t.Errorf("expected password to be 12345 but %v", *p)
	}
}
func TestWithConfigJSONMapValue(t *testing.T) {
	fs := NewFlagSet("example", ExitOnError)
	err := fs.LoadCfg("./testing/demo.json")
	if err != nil {
		t.Fatal(err)
	}
	p := fs.String("password", "", "password", BindFlag("p"), BindCfg("database"))
	err = fs.Parse(make([]string, 0))
	if err != nil {
		t.Fatal(err)
	}
	if *p != `{"password":"12345"}` {
		t.Errorf("expected password to be 12345 but %v", *p)
	}
}
func TestWithConfigYMLMapValue(t *testing.T) {
	fs := NewFlagSet("example", ExitOnError)
	err := fs.LoadCfg("./testing/demo.yaml")
	if err != nil {
		t.Fatal(err)
	}
	p := fs.String("password", "", "password", BindFlag("p"), BindCfg("database"))
	err = fs.Parse(make([]string, 0))
	if err != nil {
		t.Fatal(err)
	}
	if *p != `{"password":"12345"}` {
		t.Errorf("expected password to be 12345 but %v", *p)
	}
}

func TestBindEnv(t *testing.T) {
	os.Setenv("PASSWORD", "xyz")

	fs := NewFlagSet("example", ExitOnError)
	help := fs.String("password", "", "password", BindFlag("p"), BindEnv("PASSWORD"))
	err := fs.Parse(os.Args)
	if err != nil {
		t.Fatal(err)
	}
	f1 := fs.Lookup("password")
	f2 := fs.Lookup("p")
	if *help != "xyz" || f1 == nil || f2 == nil {
		t.Errorf("expected password to be xyz but %v", *help)
	}
}
func TestBindEnv2(t *testing.T) {
	os.Setenv("PSSWD", "xyz")

	fs := NewFlagSet("example", ExitOnError)
	help := fs.String("password", "", "password", BindFlag("p"), BindEnv("PSSWD", "PASSWORD"))
	err := fs.Parse(os.Args)
	if err != nil {
		t.Fatal(err)
	}
	f1 := fs.Lookup("password")
	f2 := fs.Lookup("p")
	if *help != "xyz" || f1 == nil || f2 == nil {
		t.Errorf("expected password to be xyz but %v", *help)
	}
}

func TestBindFlag(t *testing.T) {
	var argss = [][]string{
		{"-help"},
		{"-h"},
	}
	for _, args := range argss {
		os.Args = args

		fs := NewFlagSet("example", ExitOnError)
		/* rs := fs.EnumString("sFlag", "yes", "enter the val for s", "yes", "no", "maybe")
		rs2 := fs.EnumInt("iFlag", 1, "enter the val for i", 1, 2)
		rs3 := fs.EnumFloat("fFlag", 1.1, "enter the val for f", 1.1, 2.3) */
		help := fs.Bool("help", false, "prints help for this program", BindFlag("h"))
		err := fs.Parse(os.Args)
		if err != nil {
			t.Errorf("error not expected :%v", err)
		}
		f1 := fs.Lookup("help")
		f2 := fs.Lookup("h")
		if !*help || f1 == nil || f2 == nil {
			t.Error("Fataleerrd, help flag of name h and help shuld exist")
		}
	}
}
func TestEnums(t *testing.T) {
	os.Args = []string{"-s_flag", "value"}
	fs := NewFlagSet("test", ContinueOnError)

	fs.String("s_flag", "", "prints help for this program", EnumsOld("value,value2"))
	err := fs.Parse(os.Args)
	if err == nil {
		t.Errorf("expected error")
	}
	fmt.Println(err)
}
func TestEnums2(t *testing.T) {
	os.Args = []string{"-s_flag", "d"}
	fs := NewFlagSet("test", ContinueOnError)
	s_flag := ""
	fs.StringVar(&s_flag, "s_flag", "value2", "prints help for this program", EnumsOld("value", "value2"))
	err := fs.Parse(os.Args)
	if err == nil {
		t.Errorf("expected error")
	}
	if s_flag != "value2" {
		t.Error("flag value shouldnt be updated")
	}
}
func TestEnums3(t *testing.T) {
	os.Args = []string{"-s_flag", ""}
	fs := NewFlagSet("test", ContinueOnError)
	fs.String("s_flag", "", "prints help for this program", EnumsOld("value,value2"))
	err := fs.Parse(make([]string, 0))
	if err == nil {
		t.Errorf("expected error")
	}
	fmt.Println(err)
}
