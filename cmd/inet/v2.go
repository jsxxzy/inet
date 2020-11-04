// Author: d1y<chenhonzhou@gmail.com>
// 职院校园网客户端

package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

	var username string = ""
	var password string = ""

	easyGetLocalConfig(&username, &password)

	username1 := flag.String("username", "", "用户名")
	password1 := flag.String("password", "", "密码")

	flag.Parse()

	var checkFromData = len(*username1) >= 1 && len(*password1) >= 1
	var checkFromData1 = len(username) >= 1 && len(password) >= 1

	args := flag.Args()

	if len(args) >= 1 {
		var mode = args[0]
		switch mode {
		case "get": // 获取本地存储的账号
			var a, b string
			easyGetLocalConfig(&a, &b)
			if len(a) <= 1 || len(b) <= 1 {
				var msg = "获取失败,或未设置"
				fmt.Println(msg)
				return
			}
			fmt.Println("username:", a)
			fmt.Println("password:", b)
		case "fix": // 初始化本地存储的账号
			var msg = "初始化成功"
			if setConfigProfile("", "") != nil {
				msg = "初始化失败"
			}
			fmt.Println(msg)
		case "save": // 存储本地账号
			if checkFromData {
				setConfigProfile(*username1, *password1)
			} else {
				fmt.Println("请传递正确的账号密码")
			}
		case "check": // 检测是否登录
			var f = inet.HasLogin()
			var t = "未登录"
			if f {
				t = "已登录"
			}
			fmt.Println(t)
			break
		case "login": // 登录
			var u, p string
			var x bool = true
			if checkFromData {
				u, p = *username1, *password1
			} else {
				if checkFromData1 {
					u, p = username, password
				} else {
					x = false
				}
			}
			if !x {
				fmt.Println("请传入账号密码才能登录")
				return
			}
			info, err := inet.Login(u, p)
			if err != nil {
				fmt.Println("错误: ", err)
				return
			}
			var msg = info.GetMsg()
			fmt.Println(msg)
			break
		case "info": // 查询信息
			info, err := inet.QueryInfo()
			if err != nil {
				fmt.Println("查询信息失败")
				return
			}
			if info.Error() != nil {
				fmt.Println(info.Error())
				return
			}
			xTime, _ := strconv.Atoi(info.Time)
			fmt.Println("使用时长: ", getHumanTime(xTime))
			var flow = getHumanFlow(info.Flow)
			fmt.Println("使用流量: ", flow)
			fmt.Println("  用户id: ", info.UID)
			fmt.Println("内网地址: ", info.V4ip)
			break
		case "logout": // 退出登录
			err := inet.Logout()
			var s = "已注销"
			if err != nil {
				s = "注销错误, 未知错误"
			}
			fmt.Println(s)
			break
		default:
			fmt.Println("不存在的命令:(")
			break
		}
		return
	}

	printCmdUsage()
	flag.Usage()
	fmt.Println()

}

func exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func getHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home, nil
}

// 设置本地缓存
func setConfigProfile(u, p string) error {
	var f = getConfigFile()
	if len(u) == 0 && len(p) == 0 {
		return ioutil.WriteFile(f.Name(), []byte(""), 0777)
	}
	var parseStr = fmt.Sprintf("%v,%v", u, p)
	var b = []byte(parseStr)
	return ioutil.WriteFile(f.Name(), b, 0777)
}

// 获取配置文件, 不安全的方法, 切勿使用!!!
func getConfigFile() *os.File {
	homeDir, err := getHomeDir()
	if err != nil {
		panic(err)
	}
	var configfile = ".inetconfig"
	var file = filepath.Join(homeDir, configfile)
	if !exists(file) {
		var f, _ = os.Create(file)
		return f
	}
	var f, _ = os.Open(file)
	return f
}

// 解析配置文件
func parseConfig(b *os.File) (string, string, error) {
	var p = b.Name()
	var a, e = ioutil.ReadFile(p)
	if e != nil {
		return "", "", errors.New("get config file is error")
	}
	var c = string(a)
	var arr = strings.Split(c, ",")
	if len(arr) <= 1 {
		return "", "", errors.New("解析失败")
	}
	return arr[0], arr[1], nil
}

func easyGetLocalConfig(u, p *string) {
	var username, password, err = parseConfig(getConfigFile())
	if err != nil || len(username) <= 1 || len(password) <= 1 {
		return
	}
	*u = username
	*p = password
}

func printCmdUsage() {
	fmt.Println()
	fmt.Println("=============>")
	fmt.Println(" login: 登录")
	fmt.Println("   get: 获取保存的账号密码")
	fmt.Println("  save: 保存账号密码")
	fmt.Println("   fix: 初始化账号密码")
	fmt.Println(" check: 查询是否登录")
	fmt.Println("  info: 查询信息")
	fmt.Println("logout: 注销")
	fmt.Println("=============>")
	fmt.Println()
}
