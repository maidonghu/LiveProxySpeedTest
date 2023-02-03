package task

import (
	"LiveProxySpeedTest/internal/common"
	"LiveProxySpeedTest/internal/utils"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	bufferSize                     = 1024
	defaultTimeout                 = 10 * time.Second
	defaultDisableDownload         = false
	defaultTestNum                 = 20
	defaultMinSpeed        float64 = 0.0
)

var (
	CommonUrl = ""
	Timeout   = defaultTimeout
	Disable   = defaultDisableDownload

	TestCount = defaultTestNum
	MinSpeed  = defaultMinSpeed
)

func checkDownloadDefault() {
	if Timeout <= 0 {
		Timeout = defaultTimeout
	}
	if TestCount <= 0 {
		TestCount = defaultTestNum
	}
	if MinSpeed <= 0.0 {
		MinSpeed = defaultMinSpeed
	}
}

func TestDownloadSpeed(tableName string, ipSet utils.PingDelaySet) (speedSet utils.DownloadSpeedSet) {
	checkDownloadDefault()
	if Disable {
		return utils.DownloadSpeedSet(ipSet)
	}
	if len(ipSet) <= 0 { // IP数组长度(IP数量) 大于 0 时才会继续下载测速
		fmt.Println("\n[信息] 延迟测速结果 IP 数量为 0，跳过下载测速。")
		return
	}
	testNum := TestCount
	if len(ipSet) < TestCount || MinSpeed > 0 { // 如果IP数组长度(IP数量) 小于下载测速数量（-dn），则次数修正为IP数
		testNum = len(ipSet)
	}
	if testNum < TestCount {
		TestCount = testNum
	}

	fmt.Printf("开始下载测速（下载速度下限：%.2f MB/s，下载测速数量：%d，下载测速队列：%d）：\n", MinSpeed, TestCount, testNum)
	// 控制 下载测速进度条 与 延迟测速进度条 长度一致（强迫症）
	bar_a := len(strconv.Itoa(len(ipSet)))
	bar_b := "     "
	for i := 0; i < bar_a; i++ {
		bar_b += " "
	}
	bar := utils.NewBar(TestCount, bar_b, "")
	for i := 0; i < testNum; i++ {
		speed := downloadHandler(tableName, ipSet[i].IP)
		ipSet[i].DownloadSpeed = speed
		// 在每个 IP 下载测速后，以 [下载速度下限] 条件过滤结果
		if speed >= MinSpeed*1024*1024 {
			bar.Grow(1, "")
			speedSet = append(speedSet, ipSet[i]) // 高于下载速度下限时，添加到新数组中
			if len(speedSet) == TestCount {       // 凑够满足条件的 IP 时（下载测速数量 -dn），就跳出循环
				break
			}
		}
	}
	bar.Done()
	if len(speedSet) == 0 { // 没有符合速度限制的数据，返回所有测试数据
		speedSet = utils.DownloadSpeedSet(ipSet)
	}
	// 按速度排序
	sort.Sort(speedSet)
	return
}

// return download Speed
func downloadHandler(tableName string, ip *common.CustomIPAddr) float64 {
	ipStr := ip.IPAddr.IP.String()
	if strings.Contains(ipStr, ":") {
		ipStr = "[" + ipStr + "]"
	}
	if CommonUrl == "" {
		if tmpUrl := common.GetCommonUrl(tableName, ipStr); tmpUrl != "" {
			CommonUrl = tmpUrl
		} else {
			return 0.0
		}
	}
	url := strings.Replace(CommonUrl, "$ip$", ipStr, 1)
	start := time.Now()
	resp, err := common.LSTClient.R().Get(url)
	if err != nil {
		log.Println(err)
		return 0.0
	}
	size, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		log.Println(err)
		return 0.0
	}

	elapsed := time.Since(start).Seconds()

	if tableName == "bestv" {
		channelId := resp.Header.Get("X-TXlive-ChannelId")
		streamId := resp.Header.Get("X-TXlive-StreamId")
		if channelId == "" || streamId == "" {
			if channelId == "" && streamId == "" {
				ip.Note = "大概率串台"
			} else {
				ip.Note = "小概率串台"
			}
		}
	}

	speed := float64(size) / elapsed
	return speed
}
