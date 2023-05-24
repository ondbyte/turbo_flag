package main

import (
	"fmt"
	"os"

	"github.com/ondbyte/flag"
)

func main() {
	fs := flag.NewFlagSet("example", flag.ExitOnError)
	/* rs := fs.EnumString("sFlag", "yes", "enter the val for s", "yes", "no", "maybe")
	rs2 := fs.EnumInt("iFlag", 1, "enter the val for i", 1, 2)
	rs3 := fs.EnumFloat("fFlag", 1.1, "enter the val for f", 1.1, 2.3) */
	h := fs.Bool("h", false, "hjsdfhsb")
	err := fs.Parse(os.Args[1:])
	if err != nil || *h {
		fmt.Println(err)
		fs.Usage()
		return
	}
	fmt.Printf("%v\n", fs.Lookup("h"))
	/* fmt.Println("-----after Parse()-------")
	println(*rs)
	println(*rs2)
	println(*rs3)
	ok, err := fs.ParseCfgFile("./cfg.json")
	if err != nil {
		panic(err)
	}
	if !ok {
		_, err = fs.GenCfgFile("./cfg.json")
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("-----after ParseCfgFile()-------")
		println(*rs)
		println(*rs2)
		println(*rs3)
	}
	err = fs.PromptAll()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("-----after PromptAll()-------")
		println(*rs)
		println(*rs2)
		println(*rs3)
	} */
}
