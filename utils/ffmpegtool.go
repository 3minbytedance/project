package utils

import (
	"fmt"
	"os/exec"
)

func GetVideoFrame(fileName string) {
	// 设置视频源文件路径
	// 设置转码后文件路径
	outputFile := "./output.jpg"

	// 设置 ffmpeg 命令行参数
	args := []string{"-i", fileName, "-ss", "1", outputFile}

	// 创建 *exec.Cmd
	cmd := exec.Command("ffmpeg", args...)

	// 运行 ffmpeg 命令
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("转码成功")
}
