package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
)

type Flags struct {
	flagSl  bool
	flagF   bool
	flagD   bool
	flagExt string
}

func (flags *Flags) flagParser() {
	flag.BoolVar(&flags.flagSl, "sl", false, "./myFind -sl /path/to/dir")
	flag.BoolVar(&flags.flagF, "f", false, "./myFind -sl /path/to/dir")
	flag.BoolVar(&flags.flagD, "d", false, "./myFind -sl /path/to/dir")
	flag.StringVar(&flags.flagExt, "ext", "", "./myFind -f -ext '.extension' /path/to/dir")
	flag.Parse()
	if flag.NArg() != 1 || (!flags.flagF && flags.flagExt != "") {
		flag.PrintDefaults()
		log.Fatal("usage error")
	}
	if !flags.flagF && !flags.flagD && !flags.flagSl {
		flags.flagF = true
		flags.flagD = true
		flags.flagSl = true
	}
}

func main() {
	userFlags := Flags{}
	userFlags.flagParser()
	filepath.Walk(flag.Arg(0), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(info.Mode())
		if userFlags.flagSl && info.Mode()&fs.ModeSymlink != 0 {
			fmt.Println(path)
		} else if userFlags.flagD && info.IsDir() {
			fmt.Println(path)
		} else if userFlags.flagF && info.Mode()>>9 == 0 {
			if userFlags.flagExt != "" {
				if filepath.Ext(path) == "."+userFlags.flagExt {
					fmt.Println(path)
				}
			} else {
				fmt.Println(path)
			}
		}
		return nil
	})
}
