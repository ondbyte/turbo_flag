package flag

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_getFirstSubCommandWithArgs(t *testing.T) {
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
			got, got1, got2 := getFirstSubCommandWithArgs(tt.args.args)
			if got != tt.want {
				t.Errorf("getFirstSubCommandWithArgs() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getFirstSubCommandWithArgs() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("getFirstSubCommandWithArgs() got2 = %v, want %v", got2, tt.want2)
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
				fs.SubCmd("yadu", func(args []string) {})
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
	fs.SubCmd("yadu", func(args []string) {
		arg = args
		fs2 := NewFlagSet("second", ContinueOnError)
		fs2.SubCmd("nandan", func(args []string) {
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

func Test_resolvePointer(t *testing.T) {
	type args struct {
		ptr interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			"one level ptr",
			args{
				ptr: func() interface{} {
					v := 1
					return &v
				}(),
			},
			1,
		},
		{
			"two level ptr",
			args{
				ptr: func() interface{} {
					v := true
					v2 := &v
					return &v2
				}(),
			},
			true,
		},
		{
			"two level ptr",
			args{
				ptr: func() interface{} {
					v := false
					v2 := &v
					return &v2
				}(),
			},
			false,
		},
		{
			"two level ptr",
			args{
				ptr: func() interface{} {
					var v *FlagSet
					v2 := &v
					return &v2
				}(),
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolvePointer(tt.args.ptr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolvePointer() = %v, want %v", got, tt.want)
			}
		})
	}
}
