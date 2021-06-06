package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var path string
var timeRegexp = regexp.MustCompile(`^time=\d{2}:\d{2}:\d{2}.\d{2}`)

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
	scanner.Split(bufio.ScanWords)
	//scanner.Split(bufio.ScanLines)
	go func() {
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Printf("%s\n", m)
		}
	}()

	scannerErr := bufio.NewScanner(stdErr)
	scannerErr.Split(bufio.ScanWords)
	go func() {

		getOriginDurationFlag := true
		mp3FlagStr := "'output.mp3':"
		mp3Flag := false
		durationFlagStr := "Duration:"
		durationFlag := false
		totalSeconds := int64(0)
		currentSeconds := int64(0)
		startTime := int64(0)

		for scannerErr.Scan() {
			m := scannerErr.Text()
			if getOriginDurationFlag {
				if !mp3Flag && m == mp3FlagStr {
					mp3Flag = true
				}
				//
				if mp3Flag && !durationFlag && m == durationFlagStr {
					durationFlag = true
					continue
				}
				if mp3Flag && durationFlag {
					totalSeconds, err = ParseMp3Duration(m)
					if err != nil {
						println(m)
						fmt.Printf("解析视频总时长出错，请联系管理员: %s\n", m)
						time.Sleep(30 * time.Minute)
					}
					fmt.Printf("视频总时长: %s\n", m)
					getOriginDurationFlag = false
				}
			}

			if timeRegexp.MatchString(m) {
				if "time=00:00:00.00" == m {
					println("正在计算所需时间，请稍等")
					continue
				}
				if startTime == 0 {
					startTime = time.Now().Unix()
				}

				currentSeconds, err = ParseMp4Time(m)
				if err != nil {
					println(m)
					fmt.Printf("解析视频总时长出错，请联系管理员: %s\n", m)
					time.Sleep(30 * time.Minute)
				}

				progressFloat := float64(currentSeconds) / float64(totalSeconds)
				if currentSeconds == 0 {
					fmt.Printf("视频总时长: %s -- 进度: %.2f %%\n", SecondsToStr(float64(totalSeconds)), progressFloat)
				} else {
					usedTime := float64(time.Now().Unix() - startTime)
					factor := usedTime / progressFloat
					needTime := factor * (1 - progressFloat)

					fmt.Printf("视频总时长: %s -- 进度: %.2f %% -- 输出剩余时间: %s -- 已使用时间: %s \n",
						SecondsToStr(float64(totalSeconds)),
						progressFloat*100,
						SecondsToStr(needTime),
						SecondsToStr(usedTime))
				}

			}
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

	println("成功: 生成 export.mp4")
	println("可以关闭此程序，或一段时间后程序自动退出")
	time.Sleep(30 * time.Minute)
}

func SecondsToStr(t float64) string {
	ts := int64(t)
	return fmt.Sprintf("%d小时%d分钟%d秒", ts/3600, int64(math.Mod(t, 3600))/60, int64(math.Mod(t, 60)))
}

func ParseMp3Duration(s string) (int64, error) {
	s = strings.TrimSpace(s)
	layout := "15:04:05.00,"
	t, err := time.Parse(layout, s)
	if err != nil {
		return 0, err
	}
	tT, _ := time.Parse("15:04:05.00,", "00:00:00.00,")
	return int64(t.Sub(tT).Seconds()), err
}

func ParseMp4Time(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if len(s) < 16 {
		return 0, errors.New("长度错误:" + s)
	}
	s = s[5:16]
	layout := "15:04:05.00"
	t, err := time.Parse(layout, s)
	if err != nil {
		return 0, err
	}
	tT, _ := time.Parse("15:04:05.00", "00:00:00.00")
	return int64(t.Sub(tT).Seconds()), err
}
