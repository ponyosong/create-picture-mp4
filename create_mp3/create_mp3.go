package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	path              string
	listFilePath      string
	durationRegexp    = regexp.MustCompile(`Duration: \d{2}:\d{2}:\d{2}.\d{2}`)
	secondFloatRegexp = regexp.MustCompile(`\d+\.\d+`)
	fileNameRegexp    = regexp.MustCompile(`file '(.+)'`)
)

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

func init() {
	var err error

	chdir()

	path, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	println(path)

	//listFilePath = path + "\\list.txt"
	listFilePath = path + "/list.txt"
	fileNotExistPanic(listFilePath)

	//fileExistDelete(path + "\\output.mp3")
	fileExistDelete(path + "/output.mp3")
	//fileExistDelete(path + "\\timePoint.txt")
	fileExistDelete(path + "/timePoint.txt")
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

func createFileForFfmpeg(content string) {
	fileExistDelete("temp.txt")
	ffmpegFile, err := os.Create("temp.txt")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := ffmpegFile.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := ffmpegFile.WriteString("file '" + content + "'\n"); err != nil {
		panic(err)
	}
}

func main() {
	if genMp3() {
		return
	}

	if genTimePoint() {
		return
	}
	//println(utf8.DecodeRuneInString(string(out)))
	//fmt.Printf("The date is %s\n", string(out))

	println("成功: 生成 output.mp3")
	println("请生成 1.jpg 的封面后，执行: 3_make_mp4.exe")
	time.Sleep(15 * time.Second)
}

func genMp3() bool {
	var err error

	// windows
	//args := "/C "
	args := ""
	args += "-f concat -safe 0 -i list.txt -c copy output.mp3"
	//args += fmt.Sprintf("ffmpeg -f concat -safe 0 -i list.txt -c copy output.mp3")
	//args += fmt.Sprintf("ffprobe %s\\output.mp3", dir)
	//args += fmt.Sprintf("ffprobe output.mp3")
	cmd := exec.Command("ffmpeg", strings.Split(args, " ")...)

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
		return true
	}
	return false
}

func genTimePoint() bool {
	var totalSeconds int64

	file, err := os.Open("list.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	timePointFile, err := os.Create("timePoint.txt")
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := timePointFile.Close(); err != nil {
			panic(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		//fmt.Println(m)
		if fileNameRegexp.MatchString(m) {
			name, err := ParseFileName(m)
			//fmt.Println(name)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
				continue
			}

			intTime, err := getDurationByFile(name)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
				continue
			}
			//println(totalSeconds)
			text := MarshalTextSecond(totalSeconds)
			//fmt.Println(text)
			totalSeconds += intTime
			if text == "00:00:00" {
				text = "00:00:01"
			}
			s := fmt.Sprintf("%s %s\n", text, name)
			if _, err := timePointFile.WriteString(s); err != nil {
				panic(err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

func getDurationByFile(mp3FileName string) (int64, error) {
	var err error
	var currentSeconds int64

	//createFileForFfmpeg(mp3FileName)

	// windows
	//args := "/C "
	//args := ""
	////args += fmt.Sprintf(`""ffprobe" "%s""`, mp3FileName)
	//args += fmt.Sprintf(`ffprobe -v error -select_streams a:0 -show_entries stream=duration -of default=noprint_wrappers=1:nokey=1 "%s"`, mp3FileName)
	////args += `""ffprobe" "temp.txt""`
	//println(args)
	//
	//cmd := exec.Command("cmd", strings.Split(args, " ")...)
	//cmd := exec.Command("ffprobe", mp3FileName)

	s := []string{"-v", "error", "-select_streams", "a:0", "-show_entries", "stream=duration", "-of", "default=noprint_wrappers=1:nokey=1", mp3FileName}
	cmd := exec.Command("ffprobe", s...)

	stdErr, _ := cmd.StderrPipe()
	stdOut, _ := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(stdOut)
	scanner.Split(bufio.ScanLines)
	go func() {
		for scanner.Scan() {
			m := scanner.Text()
			//fmt.Printf("%s\n", m)

			if secondFloatRegexp.MatchString(m) {

				ff, err := strconv.ParseFloat(m, 32)
				if err != nil {
					fmt.Printf("解析视频总时长出错，请联系管理员: %d\n", currentSeconds)
					time.Sleep(30 * time.Minute)
				}

				currentSeconds = int64(ff)
				//currentSeconds, err = ParseMp3Time(m)
				//if err != nil {
				//	//println(m)
				//	fmt.Printf("解析视频总时长出错，请联系管理员: %d\n", currentSeconds)
				//	time.Sleep(30 * time.Minute)
				//}
			}
		}
	}()

	scannerErr := bufio.NewScanner(stdErr)
	scannerErr.Split(bufio.ScanLines)
	go func() {
		for scannerErr.Scan() {
			m := scannerErr.Text()
			//fmt.Printf("%s\n", m)
			if secondFloatRegexp.MatchString(m) {

				ff, err := strconv.ParseFloat(m, 32)
				if err != nil {
					fmt.Printf("解析视频总时长出错，请联系管理员: %d\n", currentSeconds)
					time.Sleep(30 * time.Minute)
				}

				currentSeconds = int64(ff)
				//currentSeconds, err = ParseMp3Time(m)
				//if err != nil {
				//	//println(m)
				//	fmt.Printf("解析视频总时长出错，请联系管理员: %d\n", currentSeconds)
				//	time.Sleep(30 * time.Minute)
				//}
			}
		}
	}()

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		println("失败: 获取时长")
		time.Sleep(15 * time.Second)
		return 0, err
	}

	return currentSeconds, nil
}

func MarshalTextSecond(dur int64) string {
	h := dur / (3600 * 60)
	m := dur % (3600 * 60) / (60)
	s := dur % 60

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
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

func ParseFileName(s string) (string, error) {
	s = strings.TrimSpace(s)
	if len(s) < 7 {
		return "", errors.New("长度错误:" + s)
	}
	s = s[6:]
	s = s[:len(s)-1]

	return s, nil
}
