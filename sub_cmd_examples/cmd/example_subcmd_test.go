package flag_test

import (
	"os"
	"testing"

	flag "github.com/ondbyte/turbo_flag"
)

var (
	branchName = ""
	remoteName = ""
)

func TestSubCmd(t *testing.T) {
	main()
}

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
