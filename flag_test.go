package flag_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	. "github.com/ondbyte/turbo_flag"
)

func Test_GetFirstSubCommandWithArgs(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 []string
		want2 bool
	}{
		{
			name: "sub command exists",
			args: args{
				args: []string{"yadu", "turbo", "--yes"},
			},
			want:  "yadu",
			want1: []string{"turbo", "--yes"},
			want2: true,
		},
		{
			name: "sub command doesnt exist",
			args: args{
				args: []string{"--yadu", "--turbo", "--yes"},
			},
			want:  "",
			want1: nil,
			want2: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := GetFirstSubCommandWithArgs(tt.args.args)
			if got != tt.want {
				t.Errorf("GetFirstSubCommandWithArgs() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetFirstSubCommandWithArgs() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("GetFirstSubCommandWithArgs() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestFlagSet_Parse(t *testing.T) {
	type args struct {
		arguments []string
	}
	tests := []struct {
		name    string
		fs      *FlagSet
		args    args
		wantErr error
	}{
		{
			name: "first",
			fs: func() *FlagSet {
				fs := NewFlagSet("first", ContinueOnError)
				fs.SubCmd("yadu", func(fs *FlagSet, args []string) {})
				fs.String("yadu", "", "")
				fs.String("yes", "", "")
				return fs
			}(),
			args: args{
				arguments: []string{"--yadu", "--yes"},
			},
			wantErr: nil,
		},
		{
			name: "second",
			fs: func() *FlagSet {
				fs := NewFlagSet("first", ContinueOnError)
				fs.String("yadu", "", "")
				fs.String("yes", "", "")
				return fs
			}(),
			args: args{
				arguments: []string{"yadu", "--yes"},
			},
			wantErr: fmt.Errorf("you are trying to run subcommand with name %v but it doesn't exist", "yadu"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fs
			if err := f.Parse(tt.args.arguments); tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("FlagSet.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNestedSubCMDRunsWithValidArgs(t *testing.T) {
	arg := make([]string, 0)
	arg2 := make([]string, 0)
	fs := NewFlagSet("first", ContinueOnError)
	fs.SubCmd("yadu", func(fs2 *FlagSet, args []string) {
		arg = args
		fs2.SubCmd("nandan", func(fs *FlagSet, args []string) {
			arg2 = args
		})

		fs2.String("yadu", "", "")
		fs2.String("yes", "", "")
		err := fs2.Parse(args)
		if err != nil {
			t.Fatal(err)
		}
	})
	fs.String("yadu", "", "")
	fs.String("yes", "", "")
	if len(arg) > 0 {
		t.Error("arg should be empty")
	}
	err := fs.Parse([]string{"yadu", "nandan", "--yes", "-no"})
	if err != nil {
		t.Fatal(err)
	}
	if arg[0] == "nandan" && arg[1] != "--yes" && arg[2] != "-no" {
		t.Error("should have --yes flag")
	}
	if arg2[0] != "--yes" && arg[1] != "-no" {
		t.Error("should have --yes flag")
	}
}

func TestFlagSet_BindEnv(t *testing.T) {
	os.Setenv("YADU", "123")
	os.Setenv("PORT", "3555")
	os.Setenv("YES", "TRUE")
	fs := NewFlagSet("test", ContinueOnError)
	name := fs.String("str", "", "", fs.Env("YADU"))
	port := fs.Int("port", 0, "", fs.Env("PORT"))
	yes := fs.Bool("yes", false, "", fs.Env("YES"))
	defer func() {
		err := recover()
		if err != nil {
			t.Fatalf("ecpected error to be nil %v", err)
		}
	}()

	if *name != "123" {
		t.Fatal("name should have the value")
	}
	if *port != 3555 {
		t.Fatal("port should have the value")
	}
	if !*yes {
		t.Fatal("yes should have the value")
	}
}

func TestFlagSet_BindCfg(t *testing.T) {
	for _, ext := range []string{"json", "yaml", "yml", "toml"} {
		defer func() {
			err := recover()
			if err != nil {
				t.Fatalf("expected error to be nil %v", err)
			}
		}()
		path := "./test_config/demo." + ext
		fs := NewFlagSet("test", ContinueOnError)
		exists, _ := fs.LoadCfg(path)
		if !exists {
			t.Fatalf("cfg [%v] file should exist", path)
		}
		password := fs.String("password", "", "", fs.Cfg("database.password"))

		if *password != "12345" {
			t.Fatal("expected password to be 12345")
		}
	}
	for _, ext := range []string{"json22", "yam2l", "ym2l", "tom2l"} {
		fs := NewFlagSet("test", ContinueOnError)
		password := fs.String("password", "", "", fs.Cfg("database.password"))
		defer func() {
			err := recover()
			if err != nil {
				t.Fatalf("expected error to be  nil %v", err)
			}
		}()

		_, err := fs.LoadCfg("./test_config/demo." + ext)
		if err == nil {
			t.Fatal("expected err")
		}

		if *password != "" {
			t.Fatalf("expected password to be empty but %v", *password)
		}

		if *password != "" {
			t.Fatal("expected password to be empty")
		}
	}
}

func TestFlagSet_Alias(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)
	password := fs.String("password", "", "", fs.Alias("p"))
	defer func() {
		err := recover()
		if err != nil {
			t.Fatalf("ecpected error to be nil %v", err)
		}
	}()

	err := fs.Parse([]string{"-p", "12345"})
	if err != nil {
		t.Fatal("expected no error")
	}
	if *password != "12345" {
		t.Fatal("expected password to be 12345")
	}
}

func TestFlagSet_Enums(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)
	password := fs.String("password", "12345", "", fs.Enum("12345", "123456789"))
	defer func() {
		err := recover()
		if err != nil {
			t.Fatalf("ecpected error to be nil %v", err)
		}
	}()

	options := fs.String("options", "c", "", fs.Enum("a", "b", "c"), fs.Alias("o"), fs.Env("OPTIONS", "OPTIONS2"))

	err := fs.Parse([]string{"-options", "z"})
	if err == nil {
		t.Fatal("expected error")
	}
	err = fs.Parse([]string{"-options", "a"})
	if err != nil {
		t.Fatal("expected no error")
	}
	if *password != "12345" {
		t.Fatal("expected password to be 12345")
	}
	if *options != "a" {
		t.Fatal("expected option to be a")
	}
}

func boolString(s string) string {
	if s == "0" {
		return "false"
	}
	return "true"
}

func TestEverything(t *testing.T) {
	ResetForTesting()
	Bool("test_bool", false, "bool value")
	Int("test_int", 0, "int value")
	Int64("test_int64", 0, "int64 value")
	Uint("test_uint", 0, "uint value")
	Uint64("test_uint64", 0, "uint64 value")
	String("test_string", "0", "string value")
	Float64("test_float64", 0, "float64 value")
	Duration("test_duration", 0, "time.Duration value")
	Func("test_func", "func value", func(string) error { return nil })

	m := make(map[string]*Flag)
	desired := "0"
	visitor := func(f *Flag) {
		if len(f.Name) > 5 && f.Name[0:5] == "test_" {
			m[f.Name] = f
			ok := false
			switch {
			case f.Value.String() == desired:
				ok = true
			case f.Name == "test_bool" && f.Value.String() == boolString(desired):
				ok = true
			case f.Name == "test_duration" && f.Value.String() == desired+"s":
				ok = true
			case f.Name == "test_func" && f.Value.String() == "":
				ok = true
			}
			if !ok {
				t.Error("Visit: bad value", f.Value.String(), "for", f.Name)
			}
		}
	}
	VisitAll(visitor)
	if len(m) != 9 {
		t.Error("VisitAll misses some flags")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
	m = make(map[string]*Flag)
	Visit(visitor)
	if len(m) != 0 {
		t.Errorf("Visit sees unset flags")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
	// Now set all flags
	Set("test_bool", "true")
	Set("test_int", "1")
	Set("test_int64", "1")
	Set("test_uint", "1")
	Set("test_uint64", "1")
	Set("test_string", "1")
	Set("test_float64", "1")
	Set("test_duration", "1s")
	Set("test_func", "1")
	desired = "1"
	Visit(visitor)
	if len(m) != 9 {
		t.Error("Visit fails after set")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
	// Now test they're visited in sort order.
	var flagNames []string
	Visit(func(f *Flag) { flagNames = append(flagNames, f.Name) })
	if !sort.StringsAreSorted(flagNames) {
		t.Errorf("flag names not sorted: %v", flagNames)
	}
}

func TestGet(t *testing.T) {
	ResetForTesting()
	Bool("test_bool", true, "bool value")
	Int("test_int", 1, "int value")
	Int64("test_int64", 2, "int64 value")
	Uint("test_uint", 3, "uint value")
	Uint64("test_uint64", 4, "uint64 value")
	String("test_string", "5", "string value")
	Float64("test_float64", 6, "float64 value")
	Duration("test_duration", 7, "time.Duration value")

	visitor := func(f *Flag) {
		if len(f.Name) > 5 && f.Name[0:5] == "test_" {
			g, ok := f.Value.(Getter)
			if !ok {
				t.Errorf("Visit: value does not satisfy Getter: %T", f.Value)
				return
			}
			switch f.Name {
			case "test_bool":
				ok = g.Get() == true
			case "test_int":
				ok = g.Get() == int(1)
			case "test_int64":
				ok = g.Get() == int64(2)
			case "test_uint":
				ok = g.Get() == uint(3)
			case "test_uint64":
				ok = g.Get() == uint64(4)
			case "test_string":
				ok = g.Get() == "5"
			case "test_float64":
				ok = g.Get() == float64(6)
			case "test_duration":
				ok = g.Get() == time.Duration(7)
			}
			if !ok {
				t.Errorf("Visit: bad value %T(%v) for %s", g.Get(), g.Get(), f.Name)
			}
		}
	}
	VisitAll(visitor)
}

func TestUsage(t *testing.T) {
	ResetForTesting()
	err := CommandLine.Parse([]string{"-x"})
	if err == nil {
		t.Error("parse did not fail for unknown flag")
	}
	if !strings.Contains(err.Error(), "flag provided but not defined: -x") {
		t.Error("did not get Usage for unknown flag")
	}
}

func testParse(f *FlagSet, t *testing.T) {
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	boolFlag := f.Bool("bool", false, "bool value")
	bool2Flag := f.Bool("bool2", false, "bool2 value")
	intFlag := f.Int("int", 0, "int value")
	int64Flag := f.Int64("int64", 0, "int64 value")
	uintFlag := f.Uint("uint", 0, "uint value")
	uint64Flag := f.Uint64("uint64", 0, "uint64 value")
	stringFlag := f.String("string", "0", "string value")
	float64Flag := f.Float64("float64", 0, "float64 value")
	durationFlag := f.Duration("duration", 5*time.Second, "time.Duration value")
	extra := "one-extra-argument"
	args := []string{
		"-bool",
		"-bool2=true",
		"--int", "22",
		"--int64", "0x23",
		"-uint", "24",
		"--uint64", "25",
		"-string", "hello",
		"-float64", "2718e28",
		"-duration", "2m",
		extra,
	}
	if err := f.Parse(args); err != nil {
		t.Fatal(err)
	}
	if !f.Parsed() {
		t.Error("f.Parse() = false after Parse")
	}
	if *boolFlag != true {
		t.Error("bool flag should be true, is ", *boolFlag)
	}
	if *bool2Flag != true {
		t.Error("bool2 flag should be true, is ", *bool2Flag)
	}
	if *intFlag != 22 {
		t.Error("int flag should be 22, is ", *intFlag)
	}
	if *int64Flag != 0x23 {
		t.Error("int64 flag should be 0x23, is ", *int64Flag)
	}
	if *uintFlag != 24 {
		t.Error("uint flag should be 24, is ", *uintFlag)
	}
	if *uint64Flag != 25 {
		t.Error("uint64 flag should be 25, is ", *uint64Flag)
	}
	if *stringFlag != "hello" {
		t.Error("string flag should be `hello`, is ", *stringFlag)
	}
	if *float64Flag != 2718e28 {
		t.Error("float64 flag should be 2718e28, is ", *float64Flag)
	}
	if *durationFlag != 2*time.Minute {
		t.Error("duration flag should be 2m, is ", *durationFlag)
	}
	if len(f.Args()) != 1 {
		t.Error("expected one argument, got", len(f.Args()))
	} else if f.Args()[0] != extra {
		t.Errorf("expected argument %q got %q", extra, f.Args()[0])
	}
}

func TestParse(t *testing.T) {
	ResetForTesting()
	testParse(CommandLine, t)
}

func TestFlagSetParse(t *testing.T) {
	testParse(NewFlagSet("test", ContinueOnError), t)
}

// Declare a user-defined flag type.
type flagVar []string

func (f *flagVar) String() string {
	return fmt.Sprint([]string(*f))
}

func (f *flagVar) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func TestUserDefined(t *testing.T) {
	var flags FlagSet
	flags.Init("test", ContinueOnError)
	var v flagVar
	flags.Var(&v, "v", "usage")
	if err := flags.Parse([]string{"-v", "1", "-v", "2", "-v=3"}); err != nil {
		t.Error(err)
	}
	if len(v) != 3 {
		t.Fatal("expected 3 args; got ", len(v))
	}
	expect := "[1 2 3]"
	if v.String() != expect {
		t.Errorf("expected value %q got %q", expect, v.String())
	}
}

func TestUserDefinedFunc(t *testing.T) {
	var flags FlagSet
	flags.Init("test", ContinueOnError)
	var ss []string
	flags.Func("v", "usage", func(s string) error {
		ss = append(ss, s)
		return nil
	})
	if err := flags.Parse([]string{"-v", "1", "-v", "2", "-v=3"}); err != nil {
		t.Error(err)
	}
	if len(ss) != 3 {
		t.Fatal("expected 3 args; got ", len(ss))
	}
	expect := "[1 2 3]"
	if got := fmt.Sprint(ss); got != expect {
		t.Errorf("expected value %q got %q", expect, got)
	}
	// test usage
	var buf strings.Builder
	flags.SetOutput(&buf)
	err := flags.Parse([]string{"-h"})
	if usage := err.Error(); !strings.Contains(usage, "flag provided but not defined: -h") {
		t.Errorf("usage string not included: %q", usage)
	}
	// test Func error
	flags = *NewFlagSet("test", ContinueOnError)
	flags.Func("v", "usage", func(s string) error {
		return fmt.Errorf("test error")
	})
	// flag not set, so no error
	if err := flags.Parse(make([]string, 0)); err != nil {
		t.Error(err)
	}
	// flag set, expect error
	if err := flags.Parse([]string{"-v", "1"}); err == nil {
		t.Error("expected error; got none")
	} else if errMsg := err.Error(); !strings.Contains(errMsg, "test error") {
		t.Errorf(`error should contain "test error"; got %q`, errMsg)
	}
}

// Declare a user-defined boolean flag type.
type boolFlagVar struct {
	count int
}

func (b *boolFlagVar) String() string {
	return fmt.Sprintf("%d", b.count)
}

func (b *boolFlagVar) Set(value string) error {
	if value == "true" {
		b.count++
	}
	return nil
}

func (b *boolFlagVar) IsBoolFlag() bool {
	return b.count < 4
}

func TestUserDefinedBool(t *testing.T) {
	var flags FlagSet
	flags.Init("test", ContinueOnError)
	var b boolFlagVar
	var err error
	flags.Var(&b, "b", "usage")
	if err = flags.Parse([]string{"-b", "-b", "-b", "-b=true", "-b=false", "-b", "barg", "-b"}); err != nil {
		if b.count < 4 {
			t.Error(err)
		}
	}

	if b.count != 4 {
		t.Errorf("want: %d; got: %d", 4, b.count)
	}

	if err == nil {
		t.Error("expected error; got none")
	}
}

func TestSetOutput(t *testing.T) {
	var flags FlagSet
	var buf bytes.Buffer
	flags.SetOutput(&buf)
	flags.Init("test", ContinueOnError)
	flags.Parse([]string{"-unknown"})
	if out := buf.String(); !strings.Contains(out, "-unknown") {
		t.Logf("expected output mentioning unknown; got %q", out)
	}
}

// This tests that one can reset the flags. This still works but not well, and is
// superseded by FlagSet.
func TestChangingArgs(t *testing.T) {
	ResetForTesting()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "-before", "subcmd", "-after", "args"}
	before := Bool("before", false, "")
	if err := CommandLine.Parse(os.Args[1:]); err != nil {
		t.Fatal(err)
	}
	cmd := Arg(0)
	os.Args = Args()
	after := Bool("after", false, "")
	Parse()
	args := Args()

	if !*before || cmd != "subcmd" || !*after || len(args) != 1 || args[0] != "args" {
		t.Fatalf("expected true subcmd true [args] got %v %v %v %v", *before, cmd, *after, args)
	}
}

// Test that -help invokes the usage message and returns ErrHelp.
func TestHelp(t *testing.T) {
	var helpCalled = false
	fs := NewFlagSet("help test", ContinueOnError)
	var flag bool
	fs.BoolVar(&flag, "flag", false, "regular flag")
	// Regular flag invocation should work
	err := fs.Parse([]string{"-flag=true"})
	if err != nil {
		t.Fatal("expected no err")
	}
	if !flag {
		t.Error("flag was not set by -flag")
	}
	// Help flag should work as expected.
	err = fs.Parse([]string{"-help"})
	if err == nil {
		t.Fatal("error expected")
	}
	if err.Error() != fmt.Errorf("flag provided but not defined: -help").Error() {
		t.Fatal("expected ErrHelp; got ", err)
	}
	// If we define a help flag, that should override.
	var help bool
	fs.BoolVar(&help, "help", false, "help flag")
	helpCalled = false
	err = fs.Parse([]string{"-help"})
	if err != nil {
		t.Fatal("expected no error for defined -help; got ", err)
	}
	if helpCalled {
		t.Fatal("help was called; should not have been for defined help flag")
	}
}

const defaultOutput = `"  -A\tfor bootstrapping, allow 'any' type\thas no default value\n  -Alongflagname\ndisable bounds checking\thas no default value\n  -C\ta boolean defaulting to true\tdefaults to [true]\n  -D path\nset relative path for local imports\thas no default value\n  -E string\nissue 23543\tdefaults to [0]\n  -F number\na non-zero number\tdefaults to [2.7]\n  -G float\na float that defaults to zero\thas no default value\n  -M string\na multiline\n    \thelp\n    \tstring\thas no default value\n  -N int\na non-zero int\tdefaults to [27]\n  -O\ta flag\n    \tmultiline help string\tdefaults to [true]\n  -Z int\nan int that defaults to zero\thas no default value\n  -maxT timeout\nset timeout for dial\thas no default value\n"`

/* func TestPrintDefaults(t *testing.T) {
	fs := NewFlagSet("print defaults test", ContinueOnError)
	var buf bytes.Buffer
	fs.SetOutput(&buf)
	fs.Bool("A", false, "")
	fs.Bool("Alongflagname", false, "disable bounds checking")
	fs.Bool("C", true, "a boolean defaulting to true")
	fs.String("D", "", "set relative `path` for local imports")
	fs.String("E", "0", "issue 23543")
	fs.Float64("F", 2.7, "a non-zero `number`")
	fs.Float64("G", 0, "a float that defaults to zero")
	fs.String("M", "", "a multiline\nhelp\nstring")
	fs.Int("N", 27, "a non-zero int")
	fs.Bool("O", true, "a flag\nmultiline help string")
	fs.Int("Z", 0, "an int that defaults to zero")
	fs.Duration("maxT", 0, "set `timeout` for dial")
	fs.PrintDefaults()
	got := buf.String()
	if got != defaultOutput {
		t.Errorf("got %q want %q\n", got, defaultOutput)
	}
} */

// Issue 19230: validate range of Int and Uint flag values.
func TestIntFlagOverflow(t *testing.T) {
	if strconv.IntSize != 32 {
		return
	}
	ResetForTesting()
	Int("i", 0, "")
	Uint("u", 0, "")
	if err := Set("i", "2147483648"); err == nil {
		t.Error("unexpected success setting Int")
	}
	if err := Set("u", "4294967296"); err == nil {
		t.Error("unexpected success setting Uint")
	}
}

// Issue 20998: Usage should respect CommandLine.output.
func TestUsageOutput(t *testing.T) {
	ResetForTesting()
	defer func(old []string) { os.Args = old }(os.Args)
	os.Args = []string{"app", "-i=1", "-unknown"}
	err := Parse()
	const want = "flag provided but not defined: -i"
	if got := err.Error(); got != want {
		t.Errorf("output = %q; want %q", got, want)
	}
}

func TestGetters(t *testing.T) {
	expectedName := "flag set"
	expectedErrorHandling := ContinueOnError
	expectedOutput := io.Writer(os.Stderr)
	fs := NewFlagSet(expectedName, expectedErrorHandling)

	if fs.Name() != expectedName {
		t.Errorf("unexpected name: got %s, expected %s", fs.Name(), expectedName)
	}
	if fs.ErrorHandling() != expectedErrorHandling {
		t.Errorf("unexpected ErrorHandling: got %d, expected %d", fs.ErrorHandling(), expectedErrorHandling)
	}

	expectedName = "gopher"
	expectedErrorHandling = ExitOnError
	expectedOutput = os.Stdout
	fs.Init(expectedName, expectedErrorHandling)
	fs.SetOutput(expectedOutput)

	if fs.Name() != expectedName {
		t.Errorf("unexpected name: got %s, expected %s", fs.Name(), expectedName)
	}
	if fs.ErrorHandling() != expectedErrorHandling {
		t.Errorf("unexpected ErrorHandling: got %d, expected %d", fs.ErrorHandling(), expectedErrorHandling)
	}
}

func TestParseError(t *testing.T) {
	for _, typ := range []string{"bool", "int", "int64", "uint", "uint64", "float64", "duration"} {
		fs := NewFlagSet("parse error test", ContinueOnError)
		fs.SetOutput(io.Discard)
		_ = fs.Bool("bool", false, "")
		_ = fs.Int("int", 0, "")
		_ = fs.Int64("int64", 0, "")
		_ = fs.Uint("uint", 0, "")
		_ = fs.Uint64("uint64", 0, "")
		_ = fs.Float64("float64", 0, "")
		_ = fs.Duration("duration", 0, "")
		// Strings cannot give errors.
		args := []string{"-" + typ + "=x"}
		err := fs.Parse(args) // x is not a valid setting for any flag.
		if err == nil {
			t.Errorf("Parse(%q)=%v; expected parse error", args, err)
			continue
		}
		if !strings.Contains(err.Error(), "invalid") || !strings.Contains(err.Error(), "parse error") {
			t.Errorf("Parse(%q)=%v; expected parse error", args, err)
		}
	}
}

func TestRangeError(t *testing.T) {
	bad := []string{
		"-int=123456789012345678901",
		"-int64=123456789012345678901",
		"-uint=123456789012345678901",
		"-uint64=123456789012345678901",
		"-float64=1e1000",
	}
	for _, arg := range bad {
		fs := NewFlagSet("parse error test", ContinueOnError)
		fs.SetOutput(io.Discard)
		_ = fs.Int("int", 0, "")
		_ = fs.Int64("int64", 0, "")
		_ = fs.Uint("uint", 0, "")
		_ = fs.Uint64("uint64", 0, "")
		_ = fs.Float64("float64", 0, "")
		// Strings cannot give errors, and bools and durations do not return strconv.NumError.
		err := fs.Parse([]string{arg})
		if err == nil {
			t.Errorf("Parse(%q)=%v; expected range error", arg, err)
			continue
		}
		if !strings.Contains(err.Error(), "invalid") || !strings.Contains(err.Error(), "value out of range") {
			t.Errorf("Parse(%q)=%v; expected range error", arg, err)
		}
	}
}

/* func TestExitCode(t *testing.T) {
	testenv.MustHaveExec(t)

	magic := 123
	if os.Getenv("GO_CHILD_FLAG") != "" {
		fs := NewFlagSet("test", ExitOnError)
		if os.Getenv("GO_CHILD_FLAG_HANDLE") != "" {
			var b bool
			fs.BoolVar(&b, os.Getenv("GO_CHILD_FLAG_HANDLE"), false, "")
		}
		fs.Parse([]string{os.Getenv("GO_CHILD_FLAG")})
		os.Exit(magic)
	}

	tests := []struct {
		flag       string
		flagHandle string
		expectExit int
	}{
		{
			flag:       "-h",
			expectExit: 0,
		},
		{
			flag:       "-help",
			expectExit: 0,
		},
		{
			flag:       "-undefined",
			expectExit: 2,
		},
		{
			flag:       "-h",
			flagHandle: "h",
			expectExit: magic,
		},
		{
			flag:       "-help",
			flagHandle: "help",
			expectExit: magic,
		},
	}

	for _, test := range tests {
		cmd := exec.Command(os.Args[0], "-test.run=TestExitCode")
		cmd.Env = append(
			os.Environ(),
			"GO_CHILD_FLAG="+test.flag,
			"GO_CHILD_FLAG_HANDLE="+test.flagHandle,
		)
		cmd.Run()
		got := cmd.ProcessState.ExitCode()
		// ExitCode is either 0 or 1 on Plan 9.
		if runtime.GOOS == "plan9" && test.expectExit != 0 {
			test.expectExit = 1
		}
		if got != test.expectExit {
			t.Errorf("unexpected exit code for test case %+v \n: got %d, expect %d",
				test, got, test.expectExit)
		}
	}
} */

func mustPanic(t *testing.T, testName string, expected string, f func()) {
	t.Helper()
	defer func() {
		switch msg := recover().(type) {
		case nil:
			t.Errorf("%s\n: expected panic(%q), but did not panic", testName, expected)
		case string:
			if msg != expected {
				t.Errorf("%s\n: expected panic(%q), but got panic(%q)", testName, expected, msg)
			}
		default:
			t.Errorf("%s\n: expected panic(%q), but got panic(%T%v)", testName, expected, msg, msg)
		}
	}()
	f()
}

func TestInvalidFlags(t *testing.T) {
	tests := []struct {
		flag     string
		errorMsg string
	}{
		{
			flag:     "-foo",
			errorMsg: "flag \"-foo\" begins with -",
		},
		{
			flag:     "foo=bar",
			errorMsg: "flag \"foo=bar\" contains =",
		},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("FlagSet.Var(&v, %q, \"\")", test.flag)

		fs := NewFlagSet("", ContinueOnError)
		errCh := make(chan error)
		go func() {
			defer func() {
				e := recover()
				if err, ok := e.(error); ok {
					errCh <- err
				}
				errCh <- nil
			}()
			var v flagVar
			fs.Var(&v, test.flag, "")
		}()
		var err error
		err = <-errCh
		if msg := test.errorMsg; err != nil && msg != err.Error() {
			t.Errorf("%s\n: unexpected output: expected %q, but got %q", testName, msg, err)
		}
	}
}

func TestRedefinedFlags(t *testing.T) {
	tests := []struct {
		flagSetName string
		errorMsg    string
	}{
		{
			flagSetName: "",
			errorMsg:    "flag redefined: foo",
		},
		{
			flagSetName: "fs",
			errorMsg:    "fs flag redefined: foo",
		},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("flag redefined in FlagSet(%q)", test.flagSetName)

		fs := NewFlagSet(test.flagSetName, ContinueOnError)

		var v flagVar
		fs.Var(&v, "foo", "")

		errCh := make(chan error)
		go func() {
			defer func() {
				e := recover()
				if err, ok := e.(error); ok {
					errCh <- err
				}
				errCh <- nil
			}()
			fs.Var(&v, "foo", "")
		}()
		err := <-errCh
		if msg := test.errorMsg; msg != err.Error() {
			t.Errorf("%s\n: unexpected output: expected %q, bug got %q", testName, msg, err)
		}
	}
}
