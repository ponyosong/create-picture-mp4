package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var path string
var durationRegexp = regexp.MustCompile(`Duration: \d{2}:\d{2}:\d{2}.\d{2}`)

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

func getDurationByFile(mp3FileName string) (int64, error) {
	var err error
	var currentSeconds int64

	args := "/C "
	args += fmt.Sprintf("ffprobe %s", mp3FileName)
	cmd := exec.Command("cmd", args)

	stdErr, _ := cmd.StderrPipe()
	stdOut, _ := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(stdOut)
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
			//fmt.Printf("%s\n", m)
			if durationRegexp.MatchString(m) {
				currentSeconds, err = ParseMp3Time(m)
				if err != nil {
					//println(m)
					fmt.Printf("解析视频总时长出错，请联系管理员: %d\n", currentSeconds)
					time.Sleep(30 * time.Minute)
				}

				//fmt.Printf("解析视频时长: %d\n", currentSeconds)
			}
		}
	}()

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		println("失败: 生成 output.mp3")
		time.Sleep(15 * time.Second)
		return
	}

}

func MarshalText(dur int64) string {
	h := dur / (3600 * 60 * 1000)
	m := dur % (3600 * 60 * 1000) / (60 * 1000)
	s := dur % (60 * 1000) / (1000)

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func ParseMp3Time(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if len(s) < 21 {
		return 0, errors.New("长度错误:" + s)
	}
	s = s[10:21]
	layout := "15:04:05.00"
	t, err := time.Parse(layout, s)
	if err != nil {
		return 0, err
	}
	tT, _ := time.Parse("15:04:05.00", "00:00:00.00")
	return t.Sub(tT).Milliseconds(), err
}
