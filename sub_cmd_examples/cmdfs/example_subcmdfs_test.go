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
	flag.MainCmdFs("git", "a version control implemented in golang", flag.ContinueOnError, os.Args, git)
}

func git(fs *flag.FlagSet, args []string) {
	fs.SubCmdFs("commit", "commits the changes with a message", commit)
	fs.SubCmdFs("remote", "adds a remote with name", remote)
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
