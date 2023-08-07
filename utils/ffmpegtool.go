package utils

import (
	"errors"
	"log"
	"os/exec"
)

// ffmpeg参数
const (
	inputVideoPathOption = "-i"
	startTimeOption      = "-ss"
	startTime            = "1"
)

func GetVideoFrame(fileName string) {
	if fileName == "" {
		err := errors.New("路径未指定")
		log.Fatal(err)
		return
	}
	// 设置转码后文件路径
	outputFile := "./output.jpg"

	// 设置 ffmpeg 命令行参数，获取第1s的帧
	args := []string{inputVideoPathOption, fileName, startTimeOption, startTime, outputFile}

	// 创建 *exec.Cmd
	cmd := exec.Command("ffmpeg", args...)

	// 运行 ffmpeg 命令
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
