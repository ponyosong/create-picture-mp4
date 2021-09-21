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
	p += "\\39-YoutubeMp3-1"

	//libRegEx := regexp.MustCompile("^[ 纯享 ] Jackson Wang王嘉尔《安静》《梦想的声音3》.+浙江卫视官方音乐HD.mp3$")
	//libRegEx := regexp.MustCompile(【选手CUT】陈颖恩深情演绎《时间有泪》收放自如《中国新歌声2》第6期 SING!CHINA S2 EP.6 20170818 [浙江卫视官方HD]")
	libRegEx := regexp.MustCompile("^.*\\.mp3$")

	zjwsRegEx := regexp.MustCompile("^[^\\]】]+[\\]】]([^《]+)《([^》]+)》.+浙江卫视.*\\.mp3$")
	shwsRegEx := regexp.MustCompile("^《(.+)》(.+) — .+SMG上海东方卫视.*\\.mp3$")
	gsRegEx := regexp.MustCompile("^([^\\[【《]+)《([^》]+).+(歌手官方|我是歌手).*\\.mp3$")
	hnwsRegEx := regexp.MustCompile("^《[^》]+》[^-]+-(.+)《([^》]+)》.*湖南卫视官方版.*\\.mp3$")

	e = filepath.Walk(p, func(absPath string, info os.FileInfo, err error) error {

		if err == nil && !info.IsDir() && libRegEx.MatchString(info.Name()) {
			var singerName string
			var singName string

			isMatched := false

			if zjwsRegEx.MatchString(info.Name()) {
				submatch := zjwsRegEx.FindSubmatch([]byte(info.Name()))
				singerName = strings.TrimSpace(string(submatch[1]))
				singName = strings.TrimSpace(string(submatch[2]))
				//println(singName + " - " + singerName)
				isMatched = true
				//fmt.Printf("%s\n", submatch[1])
				//fmt.Printf("%s\n", submatch[2])
			} else if shwsRegEx.MatchString(info.Name()) {
				submatch := shwsRegEx.FindSubmatch([]byte(info.Name()))
				//println(len(submatch))
				singName = strings.TrimSpace(string(submatch[1]))
				singerName = strings.TrimSpace(string(submatch[2]))
				//println(singName + " - " + singerName)

				isMatched = true
				//fmt.Printf("%s\n", submatch[1])
				//fmt.Printf("%s\n", submatch[2])
			} else if gsRegEx.MatchString(info.Name()) {
				submatch := gsRegEx.FindSubmatch([]byte(info.Name()))
				//println(len(submatch))
				singerName = strings.TrimSpace(string(submatch[1]))
				singName = strings.TrimSpace(string(submatch[2]))
				//println(singName + " - " + singerName)

				isMatched = true
				//fmt.Printf("%s\n", submatch[1])
				//fmt.Printf("%s\n", submatch[2])
			} else if hnwsRegEx.MatchString(info.Name()) {
				submatch := hnwsRegEx.FindSubmatch([]byte(info.Name()))
				//println(len(submatch))
				singerName = strings.TrimSpace(string(submatch[1]))
				singName = strings.TrimSpace(string(submatch[2]))
				//println(singName + " - " + singerName)

				isMatched = true
				//fmt.Printf("%s\n", submatch[1])
				//fmt.Printf("%s\n", submatch[2])
			}

			if !isMatched {
				return nil
			}

			name := singName + " - " + singerName + ".mp3"
			println(name)
			//name = strings.ReplaceAll(name, "(更多IT教程 微信535950311)", "")
			//println(absPath)
			//println(name)
			newPath := filepath.Dir(absPath) + "\\" + name
			err := os.Rename(absPath, newPath)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
		}
		return nil
	})
	if e != nil {
		log.Fatal(e)
	}
}
