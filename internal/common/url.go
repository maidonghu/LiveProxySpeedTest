package common

import (
	"strconv"
	"strings"
	"time"
)

func GetCommonUrl(tableName, ip string) string {
	switch tableName {
	case "aishang":
		return getAiShangTsUrl(ip)
	case "bestv":
		return getBesTVTsUrl()
	default:
		return ""
	}
}

func getBesTVTsUrl() (BesTVTsUrl string) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.FixedZone("Beijing Time", int((8 * time.Hour).Seconds()))
	}
	now := time.Now().In(loc).Add(-13 * 10 * time.Second)
	return strings.Join([]string{
		"http://$ip$/liveplay-kk.rtxapp.com/live/program/live/cctv4k/15000000/",
		now.Format("2006010215"), "/",
		strconv.Itoa(int(now.Unix()))[:9],
		".ts"}, "")
}

func getAiShangTsUrl(ip string) string {
	templatePath := "http://$ip$/live.aishang.ctlcdn.com/00000110240389_1/encoder/0/"
	templateFile := "playlist.m3u8?CONTENTID=00000110240389_1&AUTHINFO=FABqh274XDn8fkurD5614t%2B1RvYajgx%2Ba3PxUJe1SMO4OjrtFitM6ZQbSJEFffaD35hOAhZdTXOrK0W8QvBRom%2BXaXZYzB%2FQfYjeYzGgKhP%2Fdo%2BXpr4quVxlkA%2BubKvbU1XwJFRgrbX%2BnTs60JauQUrav8kLj%2FPH8LxkDFpzvkq75UfeY%2FVNDZygRZLw4j%2BXtwhj%2FIuXf1hJAU0X%2BheT7g%3D%3D&USERTOKEN=eHKuwve%2F35NVIR5qsO5XsuB0O2BhR0KR"
	m3u8Url := strings.Replace(templatePath+templateFile, "$ip$", ip, 1)
	m3u8Resp, err := LSTClient.R().Get(m3u8Url)
	if err != nil {
		return ""
	}
	m3u8RespStr := m3u8Resp.String()
	tsPath := strings.TrimSpace(m3u8RespStr[strings.Index(m3u8RespStr, ",")+1 : strings.Index(m3u8RespStr, ".ts")+3])
	return strings.Join([]string{templatePath, tsPath}, "/")
}
