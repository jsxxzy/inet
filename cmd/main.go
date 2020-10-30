// Author: d1y<chenhonzhou@gmail.com>
// 职院校园网客户端
//
//	bytefmt.ByteSize(100.5*bytefmt.MEGABYTE) // "100.5M"
//	bytefmt.ByteSize(uint64(1024)) // "1K"
//
// https://github.com/cloudfoundry/bytefmt

package main

import (
	"flag"
	"fmt"
	"math"
	"strconv"

	"github.com/jsxxzy/inet"
)

var suffixes [5]string

// =======

// Round round offset
//
// https://gist.github.com/anikitenko/b41206a49727b83a530142c76b1cb82d
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

// 获取格式化好的时间
func getHumanTime(h int) string {
	if h < 60 {
		return fmt.Sprintf("%v分钟", h)
	}
	if h == 60 {
		return "1小时"
	}
	m := h % 60
	var p float64 = 60
	b := float64(h) / p
	c := math.Floor(b)
	return fmt.Sprintf("%v小时%v分钟", c, m)
}

// getHumanFlow 转换流量格式转为阳间格式
//
// https://gist.github.com/anikitenko/b41206a49727b83a530142c76b1cb82d
func getHumanFlow(f float64) string {
	size := f * 1024 * 1024 // This is in bytes
	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	base := math.Log(size) / math.Log(1024)
	getSize := Round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	getSuffix := suffixes[int(math.Floor(base))]
	var result = strconv.FormatFloat(getSize, 'f', -1, 64) + " " + string(getSuffix)
	return result
}

// =======

func main() {

	username := flag.String("username", "", "用户名")
	password := flag.String("password", "", "密码")
	flag.Parse()

	args := flag.Args()

	if len(args) >= 1 {
		var mode = args[0]
		switch mode {
		case "check":
			var f = inet.HasLogin()
			var t = "未登录"
			if f {
				t = "已登录"
			}
			fmt.Println(t)
			break
		case "login": // 查询是否登录
			var checkFromData = len(*username) >= 1 && len(*password) >= 1
			if checkFromData {
				info, err := inet.Login(*username, *password)
				if err != nil {
					fmt.Println("错误: ", err)
					return
				}
				var msg = info.GetMsg()
				fmt.Println(msg)
			} else {
				fmt.Println("请传入账号密码才能登录")
			}
			break
		case "info": // 查询信息
			info, err := inet.QueryInfo()
			if err != nil {
				fmt.Println("查询信息失败")
				return
			}
			xTime, _ := strconv.Atoi(info.Time)
			fmt.Println("使用时长: ", getHumanTime(xTime))
			var flow = getHumanFlow(info.Flow)
			fmt.Println("使用流量(mb): ", flow)
			fmt.Println("用户id: ", info.UID)
			fmt.Println("内网地址: ", info.V4ip)
			break
		case "logout": // 退出登录
			err := inet.Logout()
			var s = "已退出"
			if err != nil {
				s = "退出错误, 未知错误"
			}
			fmt.Println(s)
			break
		}
	}

}
