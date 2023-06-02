package flag_test

import (
	"testing"

	flag "github.com/ondbyte/turbo_flag"
)

var (
	branchName = ""
	remoteName = ""
)

func TestSubCmd(t *testing.T) {
	//run our program
	git()
}

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

func commit(fs *flag.FlagSet, args []string) {
	var branch string
	fs.StringVar(&branch, "branch", "", "")
	fs.Alias(&branch, "b")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	branchName = branch
}

func remote(fs *flag.FlagSet, args []string) {
	var name string
	fs.StringVar(&name, "name", "", "")
	fs.Alias(&name, "n")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	remoteName = name
}
