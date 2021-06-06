package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var path string

func init() {
	var err error
	path, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	println(path)

	fileNotExistPanic(path + "\\list.txt")

	fileExistDelete(path + "\\output.mp3")
}

func fileNotExistPanic(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		println("发生错误，文件不存在:" + filePath)
		time.Sleep(30 * time.Second)
		panic(err)
	}
}

func fileExistDelete(filePath string) {
	if _, err := os.Stat(filePath); err == nil {
		err := os.Remove(filePath)
		if err != nil {
			println("发生错误，删除文件出错:" + filePath)
			time.Sleep(30 * time.Second)
			panic(err)
		}
	}
}

func main() {
	var err error

	args := "/C "
	args += "ffmpeg -f concat -safe 0 -i list.txt -c copy output.mp3"
	//args += fmt.Sprintf("ffmpeg -f concat -safe 0 -i list.txt -c copy output.mp3")
	//args += fmt.Sprintf("ffprobe %s\\output.mp3", dir)
	//args += fmt.Sprintf("ffprobe output.mp3")
	cmd := exec.Command("cmd", strings.Split(args, " ")...)

	stdErr, _ := cmd.StderrPipe()
	stdOut, _ := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(stdOut)
	//scanner.Split(bufio.ScanWords)
	scanner.Split(bufio.ScanLines)
	go func() {
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Printf("%s\n", m)
		}
	}()

	scannerErr := bufio.NewScanner(stdErr)
	scannerErr.Split(bufio.ScanLines)
	go func() {
		for scannerErr.Scan() {
			m := scannerErr.Text()
			fmt.Printf("%s\n", m)
		}
	}()

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		println("失败: 生成 output.mp3")
		time.Sleep(15 * time.Second)
		return
	}

	//println(utf8.DecodeRuneInString(string(out)))
	//fmt.Printf("The date is %s\n", string(out))

	println("成功: 生成 output.mp3")
	println("请生成 1.jpg 的封面后，执行: 3_make_mp4.exe")
	time.Sleep(15 * time.Second)
}
