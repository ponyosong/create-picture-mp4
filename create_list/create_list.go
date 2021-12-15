package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var mp3Regexp = regexp.MustCompile(`\.mp3$`)

func chdir() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	err = os.Chdir(exPath)
	if err != nil {
		panic(err)
	}
}

func main() {
	chdir()

	listFiles, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	ffmpegFile, err := os.Create("list.txt")
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := ffmpegFile.Close(); err != nil {
			panic(err)
		}
	}()

	for _, f := range listFiles {
		if f.IsDir() {
			continue
		}
		if !mp3Regexp.MatchString(f.Name()) || f.Name() == "output.mp3" {
			continue
		}
		if _, err := ffmpegFile.WriteString("file '" + f.Name() + "'\n"); err != nil {
			panic(err)
		}
	}

	println("成功: 生成 list.txt")
	println("请执行 2_create_mp3.exe")
	time.Sleep(15 * time.Second)

}
