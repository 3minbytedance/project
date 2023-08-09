package utils

import (
	"log"
	"os/exec"
)

// ffmpeg参数
const (
	inputVideoPathOption = "-i"
	startTimeOption      = "-ss"
	startTime            = "1" // 截取第1秒的帧
)

func GetVideoFrame(videoPath string, outputPath string) {
	if videoPath == "" {
		log.Fatal("路径未指定")
		return
	}

	// 设置 ffmpeg 命令行参数，获取第1s的帧
	args := []string{inputVideoPathOption, videoPath, startTimeOption, startTime, outputPath}

	// 创建 *exec.Cmd
	cmd := exec.Command("ffmpeg", args...)

	// 运行 ffmpeg 命令
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
		return
	}
}

func Transcoding(src string, dst string, overwrite bool) {
	args := []string{inputVideoPathOption, src, "-c:v", "libx264", "-strict", "-2", dst}
	if overwrite {
		args = append([]string{"-y"}, args...)
	}
	// 创建 *exec.Cmd
	cmd := exec.Command("ffmpeg", args...)

	// 运行 ffmpeg 命令
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
