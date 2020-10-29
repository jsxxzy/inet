// Author: d1y<chenhonzhou@gmail.com>
// 职院校园网客户端

package main

import (
	"flag"
	"fmt"

	"github.com/jsxxzy/inet"
)

func main() {
	username := flag.String("username", "", "用户名")
	password := flag.String("password", "", "密码")
	flag.Parse()

	args := flag.Args()

	if len(*username) >= 1 && len(*password) >= 1 {
		if !inet.HasLogin() {
			inet.Login(*username, *password)
		}
	}

	if len(args) >= 1 {
		var mode = args[0]
		switch mode {
		case "login": // 查询是否登录
			var f = inet.HasLogin()
			var t = "未登录"
			if f {
				t = "已登录"
			}
			fmt.Println(t)
			break
		case "info": // 查询信息
			info, err := inet.QueryInfo()
			if err != nil {
				fmt.Println("查询信息失败")
				return
			}
			fmt.Println("使用时长(分钟): ", info.Time)
			fmt.Println("使用流量(mb): ", info.Flow)
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
