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

	fileNotExistPanic(path + "\\1.jpg")
	fileNotExistPanic(path + "\\output.mp3")

	fileExistDelete(path + "\\export.mp4")
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

	return
	var err error

	args := "/C "
	args += "ffmpeg -loop 1 -i 1.jpg -i output.mp3 -c:v libx264 -c:a copy -shortest export.mp4"
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
		time.Sleep(3 * time.Second)
		return
	}

	//println(utf8.DecodeRuneInString(string(out)))
	//fmt.Printf("The date is %s\n", string(out))

	println("成功: 生成 output.mp3")
	println("请执行 3_make_mp4.exe")
	time.Sleep(3 * time.Second)
}
