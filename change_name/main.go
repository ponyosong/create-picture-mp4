package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	var (
		err error
		e   error
		p   string
	)

	p, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	libRegEx := regexp.MustCompile("^\\..+$")

	e = filepath.Walk(p, func(absPath string, info os.FileInfo, err error) error {

		if err == nil && !info.IsDir() && libRegEx.MatchString(info.Name()) {
			name := info.Name()[1:]
			name = strings.ReplaceAll(name, "(更多IT教程 微信535950311)", "")
			println(absPath)
			newPath := filepath.Dir(absPath) + "\\" + name
			err := os.Rename(absPath, newPath)

			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			//println(absPath)
		}
		return nil
	})
	if e != nil {
		log.Fatal(e)
	}
}
