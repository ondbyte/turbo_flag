package flag_test

import (
	"fmt"
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
	flag.MainCmd("git", "a version control implemented in golang", flag.ContinueOnError, os.Args, git)
}

func git(cmd flag.CMD, args []string) {
	cmd.Bool("help", false, "help message")
	cmd.SubCmd("commit", "commits the changes with a message", commit)
	cmd.SubCmd("remote", "adds a remote with name", remote)
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
	help, err := cmd.GetDefaultUsage()
	if err != nil {
		panic("error: " + err.Error())
	}
	fmt.Println(help)
}

func commit(cmd flag.CMD, args []string) {
	var branch string
	cmd.StringVar(&branch, "branch", "main", "branch name to work", cmd.Cfg("branch.name"), cmd.Alias("b"), cmd.Env("BRANCH", "MAIN_BRNCH"), cmd.Enum("main", "stable"))

	err := cmd.Parse(args)
	if err != nil {
		panic(err)
	}
	branchName = branch
	help, err := cmd.GetDefaultUsage()
	if err != nil {
		panic("error: " + err.Error())
	}
	fmt.Println(help)
}

func remote(fs flag.CMD, args []string) {
	var name string
	fs.StringVar(&name, "name", "", "remote name work with", fs.Alias("n"))

	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	remoteName = name
}
