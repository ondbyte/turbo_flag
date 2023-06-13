package flag

import (
	"encoding"
	"time"
)

type CMD interface {
	BoolVar(p *bool, name string, value bool, usage string, features ...*flagFeature)
	Bool(name string, value bool, usage string, features ...*flagFeature) *bool
	IntVar(p *int, name string, value int, usage string, features ...*flagFeature)
	Int(name string, value int, usage string, features ...*flagFeature) *int
	Int64Var(p *int64, name string, value int64, usage string, features ...*flagFeature)
	Int64(name string, value int64, usage string, features ...*flagFeature) *int64
	UintVar(p *uint, name string, value uint, usage string, features ...*flagFeature)
	Uint(name string, value uint, usage string, features ...*flagFeature) *uint
	Uint64Var(p *uint64, name string, value uint64, usage string, features ...*flagFeature)
	Uint64(name string, value uint64, usage string, features ...*flagFeature) *uint64
	StringVar(p *string, name string, value string, usage string, features ...*flagFeature)
	String(name string, value string, usage string, features ...*flagFeature) *string
	Float64Var(p *float64, name string, value float64, usage string, features ...*flagFeature)
	Float64(name string, value float64, usage string, features ...*flagFeature) *float64
	DurationVar(p *time.Duration, name string, value time.Duration, usage string, features ...*flagFeature)
	Duration(name string, value time.Duration, usage string, features ...*flagFeature) *time.Duration
	TextVar(p encoding.TextUnmarshaler, name string, value encoding.TextMarshaler, usage string, features ...*flagFeature)
	Var(value Value, name string, usage string, features ...*flagFeature)
	Name() string
	Set(name, value string) error
	GetDefaultUsage() (usage string, err error)
	Args() []string
	Func(name, usage string, fn func(string) error, features ...*flagFeature)
	ParseWithoutArgs(args []string) error
	Parse(arguments []string) error
	Parsed() bool
}
