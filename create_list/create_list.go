package main

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"
)

var mp3Regexp = regexp.MustCompile(`\.mp3$`)

func main() {
	listFiles, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	//mp3File, err := os.Create("封面用.txt")
	//if err != nil {
	//	panic(err)
	//}

	ffmpegFile, err := os.Create("list.txt")
	if err != nil {
		panic(err)
	}

	defer func() {
		//if err := mp3File.Close(); err != nil {
		//	panic(err)
		//}
		if err := ffmpegFile.Close(); err != nil {
			panic(err)
		}
	}()

	for _, f := range listFiles {
		if f.IsDir() {
			continue
		}
		if !mp3Regexp.MatchString(f.Name()) {
			continue
		}
		//if _, err := mp3File.WriteString(f.Name()+"\n"); err != nil {
		//	panic(err)
		//}
		if _, err := ffmpegFile.WriteString("file '" + f.Name() + "'\n"); err != nil {
			panic(err)
		}
	}

	println("成功: 生成 list.txt")
	println("请执行 2_create_mp3.exe")
	time.Sleep(3 * time.Second)

}
