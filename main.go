package main

import (
	"LiveProxySpeedTest/internal/task"
	"LiveProxySpeedTest/internal/utils"
	"flag"
	"fmt"
	"runtime"
)

var table string

func init() {
	flag.StringVar(&table, "t", "aishang", "指定测速类型")
	flag.Parse()
}

func main() {
	task.InitRandSeed() // 置随机数种子
	// 开始延迟测速
	pingData := task.NewPing(table).Run().FilterDelay()
	// 开始下载测速
	speedData := task.TestDownloadSpeed(table, pingData)
	utils.ExportCsv(table, speedData) // 输出文件
	speedData.Print()                 // 打印结果

	if runtime.GOOS == "windows" { // 如果是 Windows 系统，则需要按下 回车键 或 Ctrl+C 退出（避免通过双击运行时，测速完毕后直接关闭）
		fmt.Printf("按下 回车键 或 Ctrl+C 退出。")
		fmt.Scanln()
	}
}
