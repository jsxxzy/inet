// create by d1y<chenhonzhou@gmail.com>
// write date 2021/03/30
// 该脚本只编译在嵌入式设备里用来不断发送心跳包

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jsxxzy/inet"
	"github.com/robfig/cron"
)

// 默认的发送心跳包时间
var timeKeepAlive = 4

type User struct {
	Username string
	Password string
}

func getUserConfig() User {
	var user = os.Getenv("dr_user")
	var psd = os.Getenv("dr_password")
	if len(user) <= 0 || len(psd) <= 0 {
		panic("额, 先配置好环境变量的用户名和密码\n用户名: dr_user\n密码: dr_password")
	}
	return User{Username: user, Password: psd}
}

func loopLogin() {
	var user = getUserConfig()
	if !inet.HasLogin() {
		log.Println("当前未登录, 尝试登录中")
		var loginInfo, errWare = inet.Login(user.Username, user.Password)
		if errWare != nil {
			log.Println("登录错误, " + errWare.Error())
		} else {
			log.Println("登录成功, " + loginInfo.GetMsg())
		}
	} else {
		log.Println("当前已登录")
	}
}

func main() {

	c := cron.New()

	var arg = os.Args
	if len(arg) >= 2 {
		var t = arg[1]
		var i, err = strconv.Atoi(t)
		if err == nil {
			timeKeepAlive = i
		}
	}

	SpaceX := fmt.Sprintf("*/%d * * * * ?", timeKeepAlive)

	c.AddFunc(SpaceX, loopLogin)

	c.Run()

}
