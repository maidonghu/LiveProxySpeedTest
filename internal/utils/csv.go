package utils

import (
	"LiveProxySpeedTest/internal/common"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	maxDelay = 9999 * time.Millisecond
	minDelay = 0 * time.Millisecond
)

var (
	InputMaxDelay = maxDelay
	InputMinDelay = minDelay
	Output        string
	PrintNum      = 20
)

type PingData struct {
	IP       *common.CustomIPAddr
	Sended   int
	Received int
	Delay    time.Duration
}

type IPData struct {
	*PingData
	recvRate      float32
	DownloadSpeed float64
}

func (ipData *IPData) getRecvRate() float32 {
	if ipData.recvRate == 0 {
		pingLost := ipData.Sended - ipData.Received
		ipData.recvRate = float32(pingLost) / float32(ipData.Sended)
	}
	return ipData.recvRate
}

func (ipData *IPData) toString() []string {
	result := make([]string, 6)
	result[0] = ipData.IP.IPAddr.String()
	result[1] = ipData.IP.Loc
	result[2] = strconv.FormatFloat(float64(ipData.getRecvRate()), 'f', 2, 32)
	result[3] = strconv.FormatFloat(ipData.Delay.Seconds()*1000, 'f', 2, 32)
	result[4] = strconv.FormatFloat(ipData.DownloadSpeed/1024/1024, 'f', 2, 32)
	result[5] = ipData.IP.Note
	return result
}

func ExportCsv(tableName string, data []IPData) {
	cwd, _ := os.Getwd()
	fileName := strings.Join([]string{tableName, "_", time.Now().Format("20060102150405"), ".csv"}, "")
	Output = strings.Join([]string{cwd, fileName}, "/")
	fp, err := os.Create(Output)
	if err != nil {
		log.Fatalf("创建文件[%s]失败：%v", Output, err)
		return
	}
	defer fp.Close()
	fp.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，防止中文乱码
	w := csv.NewWriter(fp)         // 创建一个新的写入文件流
	_ = w.Write([]string{"IP 地址", "地理位置", "丢包率", "平均延迟", "下载速度 (MB/s)", "备注"})
	_ = w.WriteAll(convertToString(data))
	w.Flush()
}

func convertToString(data []IPData) [][]string {
	result := make([][]string, 0)
	for _, v := range data {
		result = append(result, v.toString())
	}
	return result
}

type PingDelaySet []IPData

func (s PingDelaySet) FilterDelay() (data PingDelaySet) {
	if InputMaxDelay > maxDelay || InputMinDelay < minDelay {
		return s
	}
	for _, v := range s {
		if v.Delay > InputMaxDelay { // 平均延迟上限
			break
		}
		if v.Delay < InputMinDelay { // 平均延迟下限
			continue
		}
		data = append(data, v) // 延迟满足条件时，添加到新数组中
	}
	return
}

func (s PingDelaySet) Len() int {
	return len(s)
}

func (s PingDelaySet) Less(i, j int) bool {
	iRate, jRate := s[i].getRecvRate(), s[j].getRecvRate()
	if iRate != jRate {
		return iRate < jRate
	}
	return s[i].Delay < s[j].Delay
}

func (s PingDelaySet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// 下载速度排序
type DownloadSpeedSet []IPData

func (s DownloadSpeedSet) Len() int {
	return len(s)
}

func (s DownloadSpeedSet) Less(i, j int) bool {
	return s[i].DownloadSpeed > s[j].DownloadSpeed
}

func (s DownloadSpeedSet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s DownloadSpeedSet) Print() {
	if len(s) <= 0 { // IP数组长度(IP数量) 大于 0 时继续
		fmt.Println("\n[信息] 完整测速结果 IP 数量为 0，跳过输出结果。")
		return
	}
	dateString := convertToString(s) // 转为多维数组 [][]String
	if len(dateString) < PrintNum {  // 如果IP数组长度(IP数量) 小于  打印次数，则次数改为IP数量
		PrintNum = len(dateString)
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "IP 地址", "地理位置", "丢包率", "平均延迟", "下载速度 (MB/s)", "备注"})
	for i := 0; i < PrintNum; i++ {
		t.AppendRow([]interface{}{i, dateString[i][0], dateString[i][1], dateString[i][2], dateString[i][3],
			dateString[i][4], dateString[i][5]})
	}
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "#", Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Name: "IP 地址", Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Name: "地理位置", Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Name: "丢包率", Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Name: "平均延迟", Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Name: "下载速度 (MB/s)", Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Name: "备注", Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
	})
	t.SetStyle(table.StyleLight)
	t.Render()
	fmt.Printf("\n完整测速结果已写入 %v 文件，可使用记事本/表格软件查看。\n", Output)
}
